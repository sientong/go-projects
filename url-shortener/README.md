# URL Shortener

A simple URL shortener service written in Go with analytics capabilities.

## Features

- URL shortening using SHA-256 hashing
- URL redirection
- Asynchronous access logging
- Background analytics processing
- Redis-based storage

## Prerequisites

- Go 1.20 or higher
- Redis server running on localhost:6379

## System Architecture

The system consists of several components:
- Main HTTP server for URL shortening and redirection
- Asynchronous access logger using goroutines and channels
- Background worker for analyzing access patterns
- Redis for storing URL mappings and analytics data

## Installation

1. Clone the repository
2. Run `go mod download`
3. Ensure Redis is running
4. Start the server with `go run main.go`

## Usage

### Shorten URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url"}'
```

### Access shortened URL

Visit `http://localhost:8080/{shortURL}` in your browser.

## Analytics

The system automatically collects and analyzes:
- Access timestamps
- User agents
- IP addresses
- Access patterns

Analytics data is processed every 5 minutes by the background worker.

## Project Structure

```markdown
url-shortener/
├── analytics/          # Analytics and logging
├── handlers/          # HTTP handlers
├── log/              # Log files (gitignored)
├── main.go           # Entry point
└── README.md         # This file
```
