package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"websocket/chatting-app/utils"
)

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
