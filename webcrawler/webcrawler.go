package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Crawler settings
const (
	maxDepth  = 2 // Limit crawling depth
	numWorker = 5 // Number of concurrent workers
)

// Visited URLs map
var visited sync.Map

// Fetch the HTML of a URL
func fetchURL(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	links := extractLinks(resp)
	return links, nil
}

// Extract links from a webpage
func extractLinks(resp *http.Response) []string {
	var links []string
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

// Worker function to process URLs
func worker(ctx context.Context, wg *sync.WaitGroup, jobs chan struct {
	url   string
	depth int
}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			url, depth := job.url, job.depth // depth is now an integer

			if _, seen := visited.Load(url); seen || depth >= maxDepth {
				continue
			}

			fmt.Println("Crawling:", url)
			visited.Store(url, true)

			links, err := fetchURL(url)
			if err != nil {
				fmt.Println("Error fetching:", url, err)
				continue
			}

			for _, link := range links {
				go func(link string) { // Send new links to the queue
					jobs <- struct {
						url   string
						depth int
					}{link, depth + 1}
				}(link)
			}
		}
	}
}

// Web Crawler with worker pool
func crawl(startURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30s timeout
	defer cancel()

	jobs := make(chan struct {
		url   string
		depth int
	}, 100)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go worker(ctx, &wg, jobs)
	}

	// Seed the crawler
	jobs <- struct {
		url   string
		depth int
	}{startURL, 0}

	// Wait for workers to finish
	wg.Wait()
	close(jobs)
}

func main() {
	startURL := "https://github.com/sientong"
	crawl(startURL)
}
