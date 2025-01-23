package handlers

import (
	"log"
	"net/http"
	"websocket/chatting-app/utils"
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Establish connection")

	currentConn := utils.EstablishNewConnection(w, r)
	if currentConn != nil {
		go utils.HandleIO(currentConn)
	}

	log.Println("Total connection:", len(utils.WebsocketConnections))

}
