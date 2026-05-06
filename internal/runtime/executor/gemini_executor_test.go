package executor

import (
	"net/http/httptest"
	"testing"

	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/auth"
)

func TestResolveGeminiURL(t *testing.T) {
	t.Run("standard endpoint", func(t *testing.T) {
		got, custom := resolveGeminiURL(nil, "gemini-2.5-pro", "generateContent")
		want := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-pro:generateContent"
		if got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
		if custom {
			t.Fatal("expected standard endpoint")
		}
	})

	t.Run("custom endpoint with path", func(t *testing.T) {
		auth := &cliproxyauth.Auth{
			Attributes: map[string]string{
				"base_url": "https://example.com/v1/chat/completions",
			},
		}
		got, custom := resolveGeminiURL(auth, "gemini-2.5-pro", "generateContent")
		if got != "https://example.com/v1/chat/completions" {
			t.Fatalf("url = %q, want custom base URL", got)
		}
		if !custom {
			t.Fatal("expected custom endpoint mode")
		}
	})
}

func TestApplyGeminiAuthHeader(t *testing.T) {
	t.Run("standard endpoint prefers api key header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "https://example.com", nil)
		applyGeminiAuthHeader(req, false, "api-key", "bearer-token")
		if got := req.Header.Get("x-goog-api-key"); got != "api-key" {
			t.Fatalf("x-goog-api-key = %q, want %q", got, "api-key")
		}
		if got := req.Header.Get("Authorization"); got != "" {
			t.Fatalf("Authorization = %q, want empty", got)
		}
	})

	t.Run("custom endpoint uses bearer auth with api key fallback", func(t *testing.T) {
		req := httptest.NewRequest("POST", "https://example.com", nil)
		applyGeminiAuthHeader(req, true, "api-key", "")
		if got := req.Header.Get("Authorization"); got != "Bearer api-key" {
			t.Fatalf("Authorization = %q, want %q", got, "Bearer api-key")
		}
		if got := req.Header.Get("x-goog-api-key"); got != "" {
			t.Fatalf("x-goog-api-key = %q, want empty", got)
		}
	})
}
