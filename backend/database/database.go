package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const DBFileName = "team_tracker.db"

func InitDB() *sqlx.DB {
	// Get absolute path to database file
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	dbPath := filepath.Join(cwd, DBFileName)

	log.Printf("Initializing database at: %s", dbPath)

	// Remove existing database if it exists
	if _, err := os.Stat(dbPath); err == nil {
		log.Printf("Removing existing database file")
		if err := os.Remove(dbPath); err != nil {
			log.Fatalf("Failed to remove existing database: %v", err)
		}
	}

	// Open new database connection
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to database at %s", dbPath)
	return db
}
