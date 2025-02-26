package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"url-shortener/analytics"
	"url-shortener/store"

	"github.com/gorilla/mux"
)

type URLRequest struct {
	LongURL string `json:"url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}

func generateShortURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(req.LongURL, "http://") && !strings.HasPrefix(req.LongURL, "https://") {
		http.Error(w, "URL must start with http:// or https://", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(req.LongURL)
	if err := store.SaveURL(r.Context(), shortURL, req.LongURL); err != nil {
		http.Error(w, "Error saving URL", http.StatusInternalServerError)
		return
	}

	response := URLResponse{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	longURL, err := store.GetURL(r.Context(), shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Log access asynchronously
	go analytics.LogAccess(analytics.AccessLog{
		ShortURL:   shortURL,
		AccessTime: time.Now(),
		UserAgent:  r.UserAgent(),
		RemoteAddr: r.RemoteAddr,
	})

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
