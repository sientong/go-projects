package main

import (
	"context"
	"log"
	"net/http"

	"url-shortener/analytics"
	"url-shortener/handlers"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()

	// Start analytics components
	analytics.StartAccessLogger(ctx)
	analytics.StartAnalyticsWorker(ctx)

	r := mux.NewRouter()

	// Add the logging middleware
	r.Use(analytics.LoggingMiddleware)

	r.HandleFunc("/shorten", handlers.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortURL}", handlers.RedirectURL).Methods("GET")

	log.Printf("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
