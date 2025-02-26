package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Server represents a backend server
type Server struct {
	URL        string
	ActiveConn int
	mu         sync.Mutex
}

// LoadBalancer represents the load balancer
type LoadBalancer struct {
	Servers []*Server
	mu      sync.Mutex
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(servers []string) *LoadBalancer {
	lb := &LoadBalancer{}
	for _, url := range servers {
		lb.Servers = append(lb.Servers, &Server{URL: url})
	}
	return lb
}

// GetLeastConnServer returns the server with the least active connections
func (lb *LoadBalancer) GetLeastConnServer() *Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	var leastConnServer *Server
	for _, server := range lb.Servers {
		if leastConnServer == nil || server.ActiveConn < leastConnServer.ActiveConn {
			leastConnServer = server
		}
	}
	return leastConnServer
}

// ServeHTTP handles incoming HTTP requests
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.GetLeastConnServer()
	server.mu.Lock()
	server.ActiveConn++
	server.mu.Unlock()

	defer func() {
		server.mu.Lock()
		server.ActiveConn--
		server.mu.Unlock()
	}()

	resp, err := http.Get(server.URL)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write([]byte(fmt.Sprintf("Forwarded to %s", server.URL)))
}

func main() {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	lb := NewLoadBalancer(servers)

	http.HandleFunc("/", lb.ServeHTTP)
	fmt.Println("Load Balancer started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
