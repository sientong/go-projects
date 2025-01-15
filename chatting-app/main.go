package main

import (
	"fmt"
	"net/http"
)

type M map[string]interface{}

func main() {
	mux := new(CustomMux)
	mux.RegisterMiddleware(JWTAuthMiddleware)

	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/index", IndexHandler)

	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":8080"

	fmt.Println("Listening on", server.Addr)
	server.ListenAndServe()
}
