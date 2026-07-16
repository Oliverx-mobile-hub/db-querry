package api

import "testing"

func TestNaturalQueryRequestSupportsPromtAlias(t *testing.T) {
	req := naturalQueryRequest{Promt: "查询用户"}
	prompt := req.Prompt
	if prompt == "" {
		prompt = req.Promt
	}
	if prompt != "查询用户" {
		t.Fatalf("unexpected prompt: %s", prompt)
	}
}

