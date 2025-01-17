package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"websocket/chatting-app/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var data utils.RequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	log.Println("username " + data.Username + " password " + data.Password)

	ok, userInfo := utils.Authenticate(data)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err, tokenString := utils.GenerateToken(userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	w.Write([]byte("Welcome " + userInfo["username"].(string)))
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join(utils.BASE_FILE_PATH, "chat.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Chat page served using rendering")
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
