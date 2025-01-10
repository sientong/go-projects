package albums

import (
	"database/sql"
	"fmt"
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

// GetAlbums queries the database for all albums and returns the results
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

// albumsByArtist queries the database for albums by the specified artist
func AlbumByTitle(title string) (Album, error) {
	var alb Album

	// Connect to the database
	Connect()

	// Query the database for the album with the specified ID
	row := db.QueryRow("SELECT * FROM album WHERE title = $1", title)

	// Scan the row data into the Album struct
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		return alb, err
	}

	return alb, nil
}

// AddAlbum inserts a new album into the database
func AddAlbum(c *gin.Context) (int64, error) {

	var newAlbum Album

	if err := c.BindJSON(&newAlbum); err != nil {
		return 0, fmt.Errorf("failed to add new album: %v", err)
	}

	// Check if the title and artist are not empty
	if (newAlbum.Title == "") || (newAlbum.Artist == "") {
		return 0, fmt.Errorf("title and artist must not be empty")
	}

	var lastInsertId int64

	// Check if the album already exists
	existingAlbum, err := AlbumByTitle(newAlbum.Title)
	if existingAlbum.Title == newAlbum.Title && existingAlbum.Artist == newAlbum.Artist && existingAlbum.Price == newAlbum.Price && err == nil {
		return existingAlbum.ID, fmt.Errorf("album already exists: %s by %s", existingAlbum.Title, existingAlbum.Artist)
	}

	// Insert the album details into the database and return the ID of the new record
	err = db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", newAlbum.Title, newAlbum.Artist, newAlbum.Price).Scan(&lastInsertId)

	// Close the database connection
	Close()

	if err != nil {
		return 0, fmt.Errorf("error inserting album: %v", err)
	}

	return lastInsertId, nil
}
