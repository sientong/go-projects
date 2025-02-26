package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestGetLeastConnServer(t *testing.T) {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	lb := NewLoadBalancer(servers)

	// Simulate connections
	lb.Servers[0].ActiveConn = 5
	lb.Servers[1].ActiveConn = 2
	lb.Servers[2].ActiveConn = 3

	server := lb.GetLeastConnServer()
	if server.URL != "http://localhost:8082" {
		t.Errorf("Expected http://localhost:8082, got %s", server.URL)
	}
}

func TestGetLeastConnServerEmpty(t *testing.T) {
	lb := NewLoadBalancer([]string{})
	server := lb.GetLeastConnServer()
	if server != nil {
		t.Errorf("Expected nil, got %v", server)
	}
}

func TestGetLeastConnServerSameConn(t *testing.T) {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	lb := NewLoadBalancer(servers)
	server := lb.GetLeastConnServer()
	if server == nil {
		t.Errorf("Expected non-nil server")
	}
}
func TestGetLeastConnServerConcurrency(t *testing.T) {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	lb := NewLoadBalancer(servers)
	lb.Servers[0].ActiveConn = 0
	lb.Servers[1].ActiveConn = 0
	lb.Servers[2].ActiveConn = 0

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lb.GetLeastConnServer()
		}()
	}

	wg.Wait()
}

func TestServeHTTP(t *testing.T) {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	lb := NewLoadBalancer(servers)

	// Mock backend servers
	for _, server := range lb.Servers {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()
		server.URL = ts.URL
	}

	req := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()
	lb.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}

func TestServeHTTPBackendDown(t *testing.T) {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	lb := NewLoadBalancer(servers)

	// Mock backend servers that always return error
	for _, server := range lb.Servers {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		server.URL = ts.URL
		defer ts.Close()
	}

	req := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()
	lb.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status Service Unavailable, got %v", resp.Status)
	}
}

func TestServeHTTPConnectionCounter(t *testing.T) {
	servers := []string{
		"http://localhost:8081",
	}
	lb := NewLoadBalancer(servers)

	// Mock backend server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	lb.Servers[0].URL = ts.URL

	req := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()

	// Check initial ActiveConn
	if lb.Servers[0].ActiveConn != 0 {
		t.Errorf("Expected initial ActiveConn to be 0, got %d", lb.Servers[0].ActiveConn)
	}

	lb.ServeHTTP(w, req)

	// Check ActiveConn after request
	if lb.Servers[0].ActiveConn != 0 {
		t.Errorf("Expected ActiveConn to be 0 after request, got %d", lb.Servers[0].ActiveConn)
	}
}

func TestServeHTTPConcurrentRequests(t *testing.T) {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	lb := NewLoadBalancer(servers)

	// Mock backend servers
	for _, server := range lb.Servers {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond) // Simulate some work
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()
		server.URL = ts.URL
	}

	var wg sync.WaitGroup
	numRequests := 100
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "http://localhost:8080", nil)
			w := httptest.NewRecorder()
			lb.ServeHTTP(w, req)
			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK, got %v", resp.Status)
			}
		}()
	}
	wg.Wait()
}
