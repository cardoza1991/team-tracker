package main

import (
	"log"
	"os"
	"path/filepath"
	"team-tracker-backend/controllers"
	"team-tracker-backend/database"
	"team-tracker-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Print current working directory
	cwd, _ := os.Getwd()
	log.Printf("Current working directory: %s", cwd)

	// Clean up any existing database files
	for _, dbName := range []string{"team-tracker.db", "team_tracker.db"} {
		if _, err := os.Stat(dbName); err == nil {
			log.Printf("Removing existing database file: %s", dbName)
			os.Remove(dbName)
		}
	}

	// Initialize database
	db := database.InitDB()
	defer db.Close()

	// Run migrations
	log.Println("Running database migrations...")
	database.MigrateDB(db)

	// Verify KML file exists
	kmlPath := filepath.Join(cwd, "Hampton Roads Lost Sheep Fields.kml")
	if _, err := os.Stat(kmlPath); err != nil {
		log.Fatalf("KML file not found at %s: %v", kmlPath, err)
	}
	log.Printf("Found KML file at: %s", kmlPath)

	// Populate locations
	log.Println("Starting location population from KML...")
	if err := controllers.PopulateLocations(db, kmlPath); err != nil {
		log.Printf("Warning: Error populating locations: %v", err)
	}

	// Verify data was populated
	var locationCount int
	if err := db.Get(&locationCount, "SELECT COUNT(*) FROM locations"); err != nil {
		log.Printf("Error counting locations: %v", err)
	} else {
		log.Printf("Total locations in database: %d", locationCount)
	}

	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	// Add routes
	routes.SetupRoutes(router, db)

	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
