package albums

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

func GetAlbums(c *gin.Context) error {
	// Create a slice of Album struct
	var albums []Album

	// Connect to the database
	Connect()

	// Query the database for all albums
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return err
	}

	// Close the rows when the function returns
	defer rows.Close()

	// Close the database connection
	Close()

	// Iterate over the rows, using Scan to assign column data to struct fields
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return err
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	c.IndentedJSON(http.StatusOK, albums)
	return nil
}
