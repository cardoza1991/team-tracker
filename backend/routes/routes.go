package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type LocationVisit struct {
	ID         int       `json:"id" db:"id"`
	LocationID int       `json:"location_id" db:"location_id"`
	TeamID     int       `json:"team_id" db:"team_id"`
	VisitDate  time.Time `json:"visit_date" db:"visit_date"`
	IsPreached bool      `json:"is_preached" db:"is_preached"`
	Notes      string    `json:"notes" db:"notes"`
}

type LocationStatus struct {
	ID         int     `json:"id" db:"id"`
	Name       string  `json:"name" db:"name"`
	Latitude   float64 `json:"latitude" db:"latitude"`
	Longitude  float64 `json:"longitude" db:"longitude"`
	IsPreached bool    `json:"is_preached" db:"is_preached"`
	LastVisit  string  `json:"last_visit" db:"last_visit"`
	VisitCount int     `json:"visit_count" db:"visit_count"`
}

type Statistics struct {
	TotalLocations    int `json:"total_locations"`
	PreachedLocations int `json:"preached_locations"`
	ActiveTeams       int `json:"active_teams"`
	TotalVisits       int `json:"total_visits"`
}

func SetupRoutes(router *gin.Engine, db *sqlx.DB) {
	// Initialize database with new tables
	initializeTables(db)

	// Get all locations
	router.GET("/api/locations", func(c *gin.Context) {
		var locations []struct {
			ID        int     `json:"id" db:"id"`
			Name      string  `json:"name" db:"name"`
			Latitude  float64 `json:"latitude" db:"latitude"`
			Longitude float64 `json:"longitude" db:"longitude"`
		}

		err := db.Select(&locations, "SELECT id, name, latitude, longitude FROM locations")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch locations"})
			return
		}

		c.JSON(http.StatusOK, locations)
	})

	// Get available locations
	router.GET("/api/locations/available", func(c *gin.Context) {
		var locations []struct {
			ID        int     `json:"id" db:"id"`
			Name      string  `json:"name" db:"name"`
			Latitude  float64 `json:"latitude" db:"latitude"`
			Longitude float64 `json:"longitude" db:"longitude"`
		}

		// Modified query to handle is_preached correctly
		query := `
            SELECT id, name, latitude, longitude 
            FROM locations 
            WHERE is_preached = FALSE 
            ORDER BY name
        `

		if err := db.Select(&locations, query); err != nil {
			log.Printf("Error fetching available locations: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch locations"})
			return
		}

		c.JSON(http.StatusOK, locations)
	})

	// Record a visit to a location
	router.POST("/api/visits", func(c *gin.Context) {
		var visit LocationVisit
		if err := c.ShouldBindJSON(&visit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		query := `
            INSERT INTO location_visits 
            (location_id, team_id, visit_date, is_preached, notes) 
            VALUES (?, ?, ?, ?, ?)`

		result, err := db.Exec(query,
			visit.LocationID,
			visit.TeamID,
			time.Now(),
			visit.IsPreached,
			visit.Notes,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record visit"})
			return
		}

		// Update location's preached status if needed
		if visit.IsPreached {
			_, err = db.Exec(
				"UPDATE locations SET is_preached = true WHERE id = ?",
				visit.LocationID,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location status"})
				return
			}
		}

		id, _ := result.LastInsertId()
		visit.ID = int(id)
		c.JSON(http.StatusCreated, visit)
	})

	// Get visit history for a location
	router.GET("/api/locations/:id/visits", func(c *gin.Context) {
		locationID := c.Param("id")
		var visits []LocationVisit

		err := db.Select(&visits, `
            SELECT * FROM location_visits 
            WHERE location_id = ? 
            ORDER BY visit_date DESC`, locationID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visits"})
			return
		}

		c.JSON(http.StatusOK, visits)
	})

	// Get all locations with their status
	router.GET("/api/locations/status", func(c *gin.Context) {
		var locations []LocationStatus

		query := `
            SELECT 
                l.*,
                COALESCE(MAX(v.visit_date), '') as last_visit,
                COUNT(v.id) as visit_count,
                COALESCE(MAX(v.is_preached), false) as is_preached
            FROM locations l
            LEFT JOIN location_visits v ON l.id = v.location_id
            GROUP BY l.id`

		err := db.Select(&locations, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch location statuses"})
			return
		}

		c.JSON(http.StatusOK, locations)
	})

	// Get statistics
	router.GET("/api/statistics", func(c *gin.Context) {
		var stats Statistics

		// Get total locations
		err := db.Get(&stats.TotalLocations, "SELECT COUNT(*) FROM locations")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
			return
		}

		// Get preached locations
		err = db.Get(&stats.PreachedLocations,
			"SELECT COUNT(DISTINCT location_id) FROM location_visits WHERE is_preached = true")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
			return
		}

		// Get active teams (teams with visits in last 24 hours)
		err = db.Get(&stats.ActiveTeams, `
            SELECT COUNT(DISTINCT team_id) FROM location_visits 
            WHERE visit_date >= datetime('now', '-1 day')`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
			return
		}

		// Get total visits
		err = db.Get(&stats.TotalVisits, "SELECT COUNT(*) FROM location_visits")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
			return
		}

		c.JSON(http.StatusOK, stats)
	})

	// Create a new team
	router.GET("/api/teams", func(c *gin.Context) {
		var teams []struct {
			ID     int    `json:"id" db:"id"`
			Name   string `json:"name" db:"name"`
			Leader string `json:"leader" db:"leader"`
		}
		err := db.Select(&teams, "SELECT id, name, leader FROM teams")
		if err != nil {
			log.Printf("Error fetching teams: %v", err) // Add logging
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
			return
		}
		c.JSON(http.StatusOK, teams)
	})

	router.POST("/api/teams", func(c *gin.Context) {
		var team struct {
			Name   string `json:"name"`
			Leader string `json:"leader"`
		}
		if err := c.ShouldBindJSON(&team); err != nil {
			log.Printf("Error binding JSON: %v", err) // Add logging
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		result, err := db.Exec(
			"INSERT INTO teams (name, leader) VALUES (?, ?)",
			team.Name, team.Leader,
		)
		if err != nil {
			log.Printf("Error creating team: %v", err) // Add logging
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
			return
		}

		id, _ := result.LastInsertId()
		c.JSON(http.StatusCreated, gin.H{
			"id":     id,
			"name":   team.Name,
			"leader": team.Leader,
		})
	})

	// Update a team
	router.PUT("/api/teams/:id", func(c *gin.Context) {
		id := c.Param("id")
		var team struct {
			Name   string `json:"name"`
			Leader string `json:"leader"`
		}
		if err := c.ShouldBindJSON(&team); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec(
			"UPDATE teams SET name = ?, leader = ? WHERE id = ?",
			team.Name, team.Leader, id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Team updated successfully"})
	})

	// Delete a team
	router.DELETE("/api/teams/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM teams WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
	})

	// Plan visits for a team
	router.POST("/api/teams/:id/plan", func(c *gin.Context) {
		var plan struct {
			LocationIDs []int  `json:"location_ids"`
			Date        string `json:"date"` // Format: YYYY-MM-DD
		}

		if err := c.ShouldBindJSON(&plan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		teamID := c.Param("id")

		// Start transaction
		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
			return
		}

		// Insert each planned visit
		for _, locID := range plan.LocationIDs {
			_, err := tx.Exec(`
            INSERT INTO planned_visits (location_id, team_id, planned_date)
            VALUES (?, ?, ?)
        `, locID, teamID, plan.Date)

			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to plan visits"})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Visits planned successfully"})
	})

	// Get team assignments
	// Get team assignments - keep this GET endpoint
	router.GET("/api/teams/:id/assignments", func(c *gin.Context) {
		teamID := c.Param("id")
		var assignments []struct {
			ID            int        `json:"id" db:"id"`
			LocationID    int        `json:"location_id" db:"location_id"`
			LocationName  string     `json:"location_name" db:"location_name"`
			IsCompleted   bool       `json:"is_completed" db:"is_completed"`
			AssignedDate  time.Time  `json:"assigned_date" db:"assigned_date"`
			CompletedDate *time.Time `json:"completed_date" db:"completed_date"`
		}

		query := `
        SELECT 
            ta.id,
            ta.location_id,
            l.name as location_name,
            ta.is_completed,
            ta.assigned_date,
            ta.completed_date
        FROM team_assignments ta
        JOIN locations l ON ta.location_id = l.id
        WHERE ta.team_id = ?
        ORDER BY ta.is_completed, l.name
    `

		if err := db.Select(&assignments, query, teamID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
			return
		}

		c.JSON(http.StatusOK, assignments)
	})

	// Assign locations to team - change this to POST endpoint
	router.POST("/api/teams/:id/assignments", func(c *gin.Context) {
		teamID := c.Param("id")
		var request struct {
			LocationIDs []int `json:"location_ids"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
			return
		}

		for _, locationID := range request.LocationIDs {
			_, err := tx.Exec(`
            INSERT INTO team_assignments (team_id, location_id)
            VALUES (?, ?)
            ON CONFLICT(team_id, location_id) DO NOTHING
        `, teamID, locationID)

			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign locations"})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Locations assigned successfully"})
	})

	// Update assignment status
	router.PUT("/api/teams/:id/assignments/:assignmentId", func(c *gin.Context) {
		teamID := c.Param("id") // Changed from "teamId" to "id"
		assignmentID := c.Param("assignmentId")
		var request struct {
			IsCompleted bool `json:"is_completed"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var completedDate *time.Time
		if request.IsCompleted {
			now := time.Now()
			completedDate = &now
		}

		_, err := db.Exec(`
            UPDATE team_assignments 
            SET is_completed = ?, completed_date = ?
            WHERE id = ? AND team_id = ?
        `, request.IsCompleted, completedDate, assignmentID, teamID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update assignment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Assignment updated successfully"})
	})

	// Get team's planned visits
	router.GET("/api/teams/:id/planned", func(c *gin.Context) {
		teamID := c.Param("id")
		var planned []struct {
			LocationID   int    `json:"location_id" db:"location_id"`
			LocationName string `json:"location_name" db:"name"`
			PlannedDate  string `json:"planned_date" db:"planned_date"`
		}

		query := `
        SELECT l.id as location_id, l.name, pv.planned_date
        FROM planned_visits pv
        JOIN locations l ON pv.location_id = l.id
        WHERE pv.team_id = ?
        AND DATE(pv.planned_date) >= DATE('now')
        ORDER BY pv.planned_date, l.name
    `

		if err := db.Select(&planned, query, teamID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch planned visits"})
			return
		}

		c.JSON(http.StatusOK, planned)
	})

	// Add this to your routes.go
	router.GET("/api/visits/history", func(c *gin.Context) {
		var visits []struct {
			ID           int       `json:"id" db:"id"`
			VisitDate    time.Time `json:"visit_date" db:"visit_date"`
			TeamName     string    `json:"team_name" db:"team_name"`
			LocationName string    `json:"location_name" db:"location_name"`
			IsPreached   bool      `json:"is_preached" db:"is_preached"`
			Notes        string    `json:"notes" db:"notes"`
		}

		query := `
        SELECT 
            v.id,
            v.visit_date,
            t.name as team_name,
            l.name as location_name,
            v.is_preached,
            v.notes
        FROM location_visits v
        JOIN teams t ON v.team_id = t.id
        JOIN locations l ON v.location_id = l.id
        ORDER BY v.visit_date DESC
    `

		if err := db.Select(&visits, query); err != nil {
			log.Printf("Error fetching visit history: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visit history"})
			return
		}

		c.JSON(http.StatusOK, visits)
	})
}

func initializeTables(db *sqlx.DB) {
	// First create the visits table
	createVisitsTable := `
    CREATE TABLE IF NOT EXISTS location_visits (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        location_id INTEGER NOT NULL,
        team_id INTEGER NOT NULL,
        visit_date DATETIME NOT NULL,
        is_preached BOOLEAN NOT NULL DEFAULT FALSE,
        notes TEXT,
        FOREIGN KEY(location_id) REFERENCES locations(id),
        FOREIGN KEY(team_id) REFERENCES teams(id)
    );`

	_, err := db.Exec(createVisitsTable)
	if err != nil {
		log.Printf("Error creating visits table: %v", err)
	}

	// Check if is_preached column exists in locations table
	var hasIsPreached bool
	err = db.Get(&hasIsPreached, `
        SELECT COUNT(*) > 0 
        FROM pragma_table_info('locations') 
        WHERE name='is_preached'
    `)

	if err != nil {
		log.Printf("Error checking for is_preached column: %v", err)
		return
	}

	// Add the column if it doesn't exist
	if !hasIsPreached {
		_, err = db.Exec(`ALTER TABLE locations ADD COLUMN is_preached BOOLEAN DEFAULT FALSE`)
		if err != nil {
			log.Printf("Error adding is_preached column: %v", err)
		}
	}

}
