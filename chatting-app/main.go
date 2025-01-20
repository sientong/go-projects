package main

import (
	"fmt"
	"net/http"
	"websocket/chatting-app/handlers"
	"websocket/chatting-app/utils"
)

func main() {
	runMigrations()

	utils.Connect()
	defer utils.Close()

	mux := new(utils.CustomMux)
	mux.RegisterMiddleware(utils.JWTAuthMiddleware)

	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/index", handlers.IndexHandler)
	mux.HandleFunc("/chat", handlers.ChatHandler)
	mux.HandleFunc("/ws", handlers.WebSocketHandler)

	server := new(http.Server)
	server.Handler = mux
	server.Addr = utils.SERVER_HOST + ":" + utils.SERVER_PORT

	fmt.Println("Listening on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
