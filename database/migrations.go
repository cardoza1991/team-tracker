package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// MigrateDB runs database migrations
func MigrateDB(db *sqlx.DB) {
	log.Println("Starting database migrations...")

	schema := `
    CREATE TABLE IF NOT EXISTS locations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        latitude REAL NOT NULL,
        longitude REAL NOT NULL,
        is_preached BOOLEAN DEFAULT FALSE
    );

    CREATE TABLE IF NOT EXISTS teams (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        leader TEXT NOT NULL,
        location_id INTEGER,
        FOREIGN KEY(location_id) REFERENCES locations(id)
    );

    CREATE TABLE IF NOT EXISTS location_visits (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        location_id INTEGER NOT NULL,
        team_id INTEGER NOT NULL,
        visit_date DATETIME NOT NULL,
        is_preached BOOLEAN NOT NULL DEFAULT FALSE,
        notes TEXT,
        FOREIGN KEY(location_id) REFERENCES locations(id),
        FOREIGN KEY(team_id) REFERENCES teams(id)
    );

    CREATE TABLE IF NOT EXISTS planned_visits (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        location_id INTEGER NOT NULL,
        team_id INTEGER NOT NULL,
        planned_date DATE NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        status TEXT DEFAULT 'planned', -- 'planned', 'completed', 'cancelled'
        FOREIGN KEY(location_id) REFERENCES locations(id),
        FOREIGN KEY(team_id) REFERENCES teams(id),
        UNIQUE(location_id, planned_date)
    );

CREATE TABLE IF NOT EXISTS team_assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER NOT NULL,
    location_id INTEGER NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    assigned_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_date DATETIME,
    FOREIGN KEY(team_id) REFERENCES teams(id),
    FOREIGN KEY(location_id) REFERENCES locations(id),
    UNIQUE(team_id, location_id)
);`

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")
}
