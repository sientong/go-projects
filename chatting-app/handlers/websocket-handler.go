package handlers

import (
	"log"
	"net/http"
	"websocket/chatting-app/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	currentGorillaConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	username := r.URL.Query().Get("username")
	log.Println("New client connected:", username)

	currentConn := utils.WebSocketConnection{Conn: currentGorillaConn, Username: username}
	utils.WebsocketConnections = append(utils.WebsocketConnections, &currentConn)

	log.Println("Total connection:", len(utils.WebsocketConnections))

	go utils.HandleIO(&currentConn)
}
