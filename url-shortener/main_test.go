package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"url-shortener/analytics"
	"url-shortener/handlers"
)

func createShortURL(w http.ResponseWriter, r *http.Request) {
	analytics.LoggingMiddleware(http.HandlerFunc(handlers.ShortenURL)).ServeHTTP(w, r)
}

func TestCreateShortURL(t *testing.T) {
	payload := map[string]string{
		"url": "http://example.com",
	}
	jsonPayload, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/shorten", strings.NewReader(string(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createShortURL(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	expectedShortURL := "8OamqXBC"
	if responseBody["short_url"] != expectedShortURL {
		t.Fatalf("expected short URL %s, got %s", expectedShortURL, responseBody["short_url"])
	}
}

func TestCreateShortURLInvalidPayload(t *testing.T) {
	payload := "invalid payload"
	req := httptest.NewRequest("POST", "/shorten", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createShortURL(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", resp.Status)
	}
}

func TestCreateShortURLMissingURL(t *testing.T) {
	payload := map[string]string{}
	jsonPayload, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/shorten", strings.NewReader(string(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createShortURL(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", resp.Status)
	}
}
