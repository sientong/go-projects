# Concurrent Web Crawler

A simple concurrent web crawler implemented in Go that demonstrates the use of Goroutines and channels for efficient web scraping and indexing.

## Features

- Concurrent web page fetching using Goroutines
- Safe data handling with mutexes
- Content indexing of fetched pages
- Error handling for failed requests

## Prerequisites

- Go 1.11 or higher
- Internet connection to fetch web pages


## Usage

Run the crawler with:

```bash
go run webcrawler.go
```

The program will:
1. Fetch content from predefined URLs concurrently
2. Store the content in an index
3. Print the fetching results
4. Display the first 100 characters of indexed content for each URL

## Configuration

To modify the target URLs, edit the `urls` slice in the `main()` function:

```go
urls := []string{
    "https://example.com",
    "https://golang.org",
    "https://github.com",
}
```

## License

MIT License
