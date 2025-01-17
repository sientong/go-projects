package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/novalagung/gubrak/v2"
)

type SocketPayload struct {
	Message string
}

type SocketResponse struct {
	From    string
	Type    string
	Message string
}

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

var WebsocketConnections = make([]*WebSocketConnection, 0)

func HandleIO(currentCon *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error:", fmt.Sprintf("%v", r))
		}
	}()

	broadcastMessage(currentCon, MESSAGE_NEW_USER, " ")

	for {
		payload := SocketPayload{}
		err := currentCon.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				broadcastMessage(currentCon, MESSAGE_LEAVE, " ")
				ejectConnection(currentCon)
				return
			}

			log.Println("Error:", err)
			continue
		}

		broadcastMessage(currentCon, MESSAGE_CHAT, payload.Message)
	}
}

func ejectConnection(currentConn *WebSocketConnection) {
	filtered := gubrak.From(WebsocketConnections).Reject(func(each *WebSocketConnection) bool {
		return each == currentConn
	}).Result()
	WebsocketConnections = filtered.([]*WebSocketConnection)
}

func broadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	log.Println("Broadcasting message:", message, " to ", len(WebsocketConnections), " connections")

	for _, eachConn := range WebsocketConnections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})
	}
}
