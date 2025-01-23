package utils

import (
	"fmt"
	"log"
	"net/http"
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
	Username  string
	ChannelId string
}

var WebsocketConnections = make([]*WebSocketConnection, 0)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleIO(currentCon *WebSocketConnection) {
	log.Println("Receiving message ")

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

		log.Printf("Received message from %s: %s\n", currentCon.Username, payload.Message)

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

func EstablishNewConnection(w http.ResponseWriter, r *http.Request) *WebSocketConnection {

	currentGorillaConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	username := r.URL.Query().Get("username")
	channelId := r.URL.Query().Get("channelId")
	log.Println("New client connected:", username, " on channel ", channelId)

	currentConn := WebSocketConnection{Conn: currentGorillaConn, Username: username, ChannelId: channelId}

	if len(WebsocketConnections) == 0 {
		WebsocketConnections = append(WebsocketConnections, &currentConn)
	} else {
		for _, eachConn := range WebsocketConnections {
			log.Println(eachConn.ChannelId, " ", channelId, " ", eachConn.Username, " ", username)
			if eachConn.Username == username && eachConn.ChannelId == channelId {
				return nil
			}
		}
		WebsocketConnections = append(WebsocketConnections, &currentConn)
	}

	return &currentConn
}
