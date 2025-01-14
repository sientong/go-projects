package routers

import (
	"log"

	albums "example/web-service-gin/album-functions"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init(address string, port string) {

	router := gin.Default()

	router.GET("/albums", func(c *gin.Context) {
		err := albums.GetAlbums(c)
		if err != nil {
			log.Println(err)
		}
	})

	router.POST("/albums", func(c *gin.Context) {
		lastInsertId, err := albums.AddAlbum(c)
		if err != nil {
			log.Println(err)
		} else {
			c.JSON(200, gin.H{
				"Album is successfully created with id": lastInsertId})
		}
	})
	router.Run(address + ":" + port)

}
