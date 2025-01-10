package main

import (
	"fmt"
	"log"

	albums "example/data-access/functions"

	_ "github.com/lib/pq" // underscore is required to import the package but not use it directly (only use the init function of the package
)

func main() {

	albums.Connect()

	// Query the database for all albums
	allAlbums, err := albums.GetAlbums()
	if err != nil {
		log.Fatal(err)
	}

	// Print the albums
	for _, alb := range allAlbums {
		fmt.Printf("%s by %s: $%.2f\n", alb.Title, alb.Artist, alb.Price)
	}

	albumsByArtist, err := albums.AlbumsByArtist("John Coltrane")
	if err != nil {
		albums.Close()
		log.Fatal(err)
	}

	// Print the albums
	for _, alb := range albumsByArtist {
		fmt.Printf("%s by %s: $%.2f\n", alb.Title, alb.Artist, alb.Price)
	}

	// Query the database for a specific album
	albumByID, err := albums.AlbumByID(2)
	if err != nil {
		albums.Close()
		log.Fatal(err)
	}

	// Print the album
	fmt.Printf("%s by %s: $%.2f\n", albumByID.Title, albumByID.Artist, albumByID.Price)

	// Add a new album to the database
	newID, err := albums.AddAlbum(albums.Album{Title: "The Modern Sound of Betty Carter", Artist: "Betty Carter", Price: 7.99})
	if err != nil {

		if err.Error() == "album already exists: The Modern Sound of Betty Carter by Betty Carter" {

			fmt.Println(err)
			alb, err := albums.GetLatestAlbum()

			if err != nil {
				albums.Close()
				log.Fatal(err)
			}

			fmt.Println("Latest album: ", alb)

			albums.Close()
			return
		}

		log.Fatal(err)
	}

	fmt.Printf("ID of new album: %d\n", newID)

	albums.Close()
}
