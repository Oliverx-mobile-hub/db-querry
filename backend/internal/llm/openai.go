package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"db-querry/backend/internal/api"
)

type OpenAIClient struct {
	apiKey  string
	baseURL string
	model   string
	wireAPI string
	client  *http.Client
}

func NewOpenAIClient(apiKey, baseURL, model, wireAPI string) OpenAIClient {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if model == "" {
		model = "gpt-4.1-mini"
	}
	if wireAPI == "" {
		wireAPI = "chat"
	}
	return OpenAIClient{apiKey: apiKey, baseURL: baseURL, model: model, wireAPI: wireAPI, client: &http.Client{Timeout: 45 * time.Second}}
}

func (c OpenAIClient) GenerateSQL(ctx context.Context, prompt string, metadata api.MetadataDocument) (api.GeneratedSQLDraft, error) {
	if c.apiKey == "" {
		return api.GeneratedSQLDraft{}, errors.New("missing openai api key")
	}
	dialect := "PostgreSQL"
	if api.NormalizeDatabaseType(metadata.DatabaseType) == api.DatabaseTypeMySQL {
		dialect = "MySQL"
	}
	system := fmt.Sprintf("You generate %s SELECT queries only. Use %s SQL dialect. Return strict JSON with sql, explanation, referencedObjects. Do not include markdown.", dialect, dialect)
	user := buildPrompt(prompt, metadata)
	if strings.EqualFold(c.wireAPI, "responses") {
		return c.generateSQLWithResponses(ctx, system, user)
	}
	return c.generateSQLWithChat(ctx, system, user)
}

func (c OpenAIClient) generateSQLWithChat(ctx context.Context, system string, user string) (api.GeneratedSQLDraft, error) {
	reqBody := map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"temperature":     0,
		"response_format": map[string]string{"type": "json_object"},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return api.GeneratedSQLDraft{}, fmt.Errorf("openai status %d: %s", resp.StatusCode, sanitizeOpenAIError(string(body)))
	}
	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	if len(parsed.Choices) == 0 {
		return api.GeneratedSQLDraft{}, errors.New("openai returned no choices")
	}
	return parseDraft(parsed.Choices[0].Message.Content)
}

func (c OpenAIClient) generateSQLWithResponses(ctx context.Context, system string, user string) (api.GeneratedSQLDraft, error) {
	reqBody := map[string]any{
		"model": c.model,
		"input": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"text": map[string]any{
			"format": map[string]string{"type": "json_object"},
		},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/responses", bytes.NewReader(body))
	if err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return api.GeneratedSQLDraft{}, fmt.Errorf("openai status %d: %s", resp.StatusCode, sanitizeOpenAIError(string(body)))
	}
	var parsed struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	text := parsed.OutputText
	if text == "" {
		for _, item := range parsed.Output {
			for _, content := range item.Content {
				if content.Text != "" {
					text = content.Text
					break
				}
			}
			if text != "" {
				break
			}
		}
	}
	if text == "" {
		return api.GeneratedSQLDraft{}, errors.New("openai returned no response text")
	}
	return parseDraft(text)
}

func sanitizeOpenAIError(body string) string {
	body = strings.ReplaceAll(body, "\n", " ")
	body = strings.TrimSpace(body)
	body = regexp.MustCompile(`sk-[A-Za-z0-9_\-*]+`).ReplaceAllString(body, "sk-[redacted]")
	if body == "" {
		return "empty error body"
	}
	return body
}

func buildPrompt(prompt string, metadata api.MetadataDocument) string {
	metadata.DatabaseType = api.NormalizeDatabaseType(metadata.DatabaseType)
	metadataJSON, _ := json.Marshal(metadata)
	text := string(metadataJSON)
	if len(text) > 24000 {
		text = text[:24000]
	}
	return fmt.Sprintf("Database type: %s\nUser request: %s\nMetadata JSON: %s\nReturn JSON: {\"sql\":\"SELECT ...\",\"explanation\":\"...\",\"referencedObjects\":[\"schema.table\"]}", metadata.DatabaseType, prompt, text)
}

func parseDraft(content string) (api.GeneratedSQLDraft, error) {
	var raw struct {
		SQL               string   `json:"sql"`
		Explanation       string   `json:"explanation"`
		ReferencedObjects []string `json:"referencedObjects"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(content)), &raw); err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	if strings.TrimSpace(raw.SQL) == "" {
		return api.GeneratedSQLDraft{}, errors.New("empty generated sql")
	}
	return api.GeneratedSQLDraft{SQL: raw.SQL, Explanation: raw.Explanation, ReferencedObjects: raw.ReferencedObjects}, nil
}
