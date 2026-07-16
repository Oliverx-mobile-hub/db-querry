package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"db-querry/backend/internal/api"
)

type OpenAIClient struct {
	apiKey string
	client *http.Client
}

func NewOpenAIClient(apiKey string) OpenAIClient {
	return OpenAIClient{apiKey: apiKey, client: &http.Client{Timeout: 45 * time.Second}}
}

func (c OpenAIClient) GenerateSQL(ctx context.Context, prompt string, metadata api.MetadataDocument) (api.GeneratedSQLDraft, error) {
	if c.apiKey == "" {
		return api.GeneratedSQLDraft{}, errors.New("missing openai api key")
	}
	system := "You generate PostgreSQL SELECT queries only. Return strict JSON with sql, explanation, referencedObjects. Do not include markdown."
	user := buildPrompt(prompt, metadata)
	reqBody := map[string]any{
		"model": "gpt-4.1-mini",
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"temperature": 0,
		"response_format": map[string]string{"type": "json_object"},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return api.GeneratedSQLDraft{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
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
		return api.GeneratedSQLDraft{}, fmt.Errorf("openai status %d", resp.StatusCode)
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

func buildPrompt(prompt string, metadata api.MetadataDocument) string {
	metadataJSON, _ := json.Marshal(metadata)
	text := string(metadataJSON)
	if len(text) > 24000 {
		text = text[:24000]
	}
	return fmt.Sprintf("User request: %s\nMetadata JSON: %s\nReturn JSON: {\"sql\":\"SELECT ...\",\"explanation\":\"...\",\"referencedObjects\":[\"schema.table\"]}", prompt, text)
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

