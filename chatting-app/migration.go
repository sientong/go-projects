package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations() {
	// Database connection URL
	// Example for PostgreSQL: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	dbURL := "postgres://postgres:099052@localhost:5432/syaif?sslmode=disable"

	// Path to migration files
	migrationsPath := "file://db/migrations"

	// Initialize migrate instance
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully.")
}
