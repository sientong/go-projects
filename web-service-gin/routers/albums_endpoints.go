package routers

import (
	"log"

	albums "example/web-service-gin/album-functions"

	"github.com/gin-gonic/gin"
)

func GetAlbumsEndpoint(c *gin.Context) {

	router.GET("/albums", func(c *gin.Context) {
		err := albums.GetAlbums(c)
		if err != nil {
			log.Println(err)
		}
	})
}

func AddAlbumEndpoint(c *gin.Context) {

	router.POST("/albums", func(c *gin.Context) {
		lastInsertId, err := albums.AddAlbum(c)
		if err != nil {
			log.Println(err)
		} else {
			c.JSON(200, gin.H{
				"Album is successfully created with id": lastInsertId})
		}
	})
}
