package handlers

import (
	"net/http"
	"websocket/chatting-app/utils"

	"github.com/labstack/echo/v4"
)

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Read the body and assign the data
// 	var data utils.RequestData
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
// 		return
// 	}
// 	defer r.Body.Close()

// 	log.Println("username " + data.Username + " password " + data.Password)

// 	tokenString, err := utils.Authenticate(data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	w.Write([]byte(tokenString))
// }

func LoginHandler(c echo.Context) error {
	var data utils.RequestData
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read request body")
	}

	token, err := utils.Authenticate(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
