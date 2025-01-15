package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNATURE_KEY = []byte("secret")
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_ISSUER = "chatting-app"
var connections = make([]*WebSocketConnection, 0)

type CustomClaim struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	Group    string `json:"group"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ok, userInfo := authenticate(username, password)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims := CustomClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    JWT_ISSUER,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		Username: userInfo["username"].(string),
		Email:    userInfo["email"].(string),
		Group:    userInfo["group"].(string),
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, _ := json.Marshal(M{"token": signedToken})
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
	content, err := os.ReadFile("chat.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Chat page served")
	fmt.Fprintf(w, "%s", content)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	currentGorillaConn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	username := r.URL.Query().Get("username")
	log.Println("New client connected:", username)

	currentConn := WebSocketConnection{Conn: currentGorillaConn, Username: username}
	connections = append(connections, &currentConn)

	log.Println("Total connection:", len(connections))

	go handleIO(&currentConn, connections)
}
