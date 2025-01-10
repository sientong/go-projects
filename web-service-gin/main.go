package main

import (
	"log"

	albums "example/web-service-gin/album-functions"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("/albums", func(c *gin.Context) {
		err := albums.GetAlbums(c)
		if err != nil {
			log.Fatal(err)
		}
	})

	router.Run("localhost:8080")
}
