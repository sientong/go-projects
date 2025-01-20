package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // underscore is required to import the package but not use it directly (only use the init function of the package
)

var db *sql.DB

func Connect() {
	var err error

	// Connect to the database
	conn := "user=postgres dbname=syaif sslmode=disable password=099052"
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if err != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to database")
}

func Close() {
	defer db.Close()
	fmt.Println("Database connection closed")
}
