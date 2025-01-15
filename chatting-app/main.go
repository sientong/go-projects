package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := new(CustomMux)
	mux.RegisterMiddleware(JWTAuthMiddleware)

	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/index", IndexHandler)
	mux.HandleFunc("/chat", ChatHandler)
	mux.HandleFunc("/ws", WebSocketHandler)

	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":8080"

	fmt.Println("Listening on", server.Addr)
	server.ListenAndServe()
}
