package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"websocket/chatting-app/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body and assign the data
	var data utils.RequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	log.Println("username " + data.Username + " password " + data.Password)

	tokenString, err := utils.Authenticate(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write([]byte(tokenString))
}
