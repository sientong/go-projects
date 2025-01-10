package albums

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // underscore is required to import the package but not use it directly (only use the init function of the package
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func GetAlbums() ([]Album, error) {
	// Create a slice of Album struct
	var albums []Album

	// Query the database for all albums
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, err
	}

	// Close the rows when the function returns
	defer rows.Close()

	// Iterate over the rows, using Scan to assign column data to struct fields
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, err
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

// albumsByArtist queries the database for albums by the specified artist
func AlbumsByArtist(name string) ([]Album, error) {

	// Create a slice of Album struct
	var albums []Album

	// Query the database for albums by the specified artist
	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, err
	}

	// Close the rows when the function returns
	defer rows.Close()

	// Iterate over the rows, using Scan to assign column data to struct fields
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, err
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

// albumByID queries the database for an album with the specified ID
func AlbumByID(id int64) (Album, error) {
	var alb Album

	// Query the database for the album with the specified ID
	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)

	// Scan the row data into the Album struct
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		return alb, err
	}

	return alb, nil
}

func AlbumByTitle(title string) (Album, error) {
	var alb Album

	// Query the database for the album with the specified ID
	row := db.QueryRow("SELECT * FROM album WHERE title = $1", title)

	// Scan the row data into the Album struct
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		return alb, err
	}

	return alb, nil
}

func GetLatestAlbum() (Album, error) {
	var alb Album

	// Query the database for the album with the specified ID
	row := db.QueryRow("SELECT * FROM album ORDER BY id DESC LIMIT 1")

	// Scan the row data into the Album struct
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		return alb, err
	}

	return alb, nil
}

// addAlbum adds a new album to the database and returns the ID of the new record
func AddAlbum(alb Album) (int64, error) {

	// Check if the title and artist are not empty
	if (alb.Title == "") || (alb.Artist == "") {
		return 0, fmt.Errorf("title and artist must not be empty")
	}

	var lastInsertId int64
	// Check if the album already exists
	existingAlbum, err := AlbumByTitle(alb.Title)
	if existingAlbum.Title == alb.Title && existingAlbum.Artist == alb.Artist && existingAlbum.Price == alb.Price && err == nil {
		return 0, fmt.Errorf("album already exists: %s by %s", existingAlbum.Title, existingAlbum.Artist)
	}

	// Insert the album details into the database and return the ID of the new record
	err = db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", alb.Title, alb.Artist, alb.Price).Scan(&lastInsertId)

	if err != nil {
		return 0, fmt.Errorf("error inserting album: %v", err)
	}

	return lastInsertId, nil
}
