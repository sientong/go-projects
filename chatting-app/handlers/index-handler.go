package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	w.Write([]byte("Welcome " + userInfo["username"].(string)))
}
