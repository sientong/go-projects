package main

import (
	"net/http"
	"websocket/chatting-app/handlers"
	"websocket/chatting-app/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	runMigrations()

	utils.Connect()
	defer utils.Close()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/login", handlers.LoginHandler)

	mux := new(utils.CustomMux)
	mux.RegisterMiddleware(utils.JWTAuthMiddleware)

	// mux.HandleFunc("/login/", handlers.LoginHandler)
	mux.HandleFunc("/index", handlers.IndexHandler)
	mux.HandleFunc("/chat", handlers.ChatHandler)
	mux.HandleFunc("/ws", handlers.WebSocketHandler)
	mux.HandleFunc("/channel/", handlers.ChannelHandler)

	// server := new(http.Server)
	// server.Handler = mux
	// server.Addr = utils.SERVER_HOST + ":" + utils.SERVER_PORT

	// fmt.Println("Listening on", server.Addr)
	// err := server.ListenAndServe()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	e.Logger.Fatal(e.Start(utils.SERVER_HOST + ":" + utils.SERVER_PORT))

}
