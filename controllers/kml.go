package controllers

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// KML structure definitions
type KML struct {
	XMLName  xml.Name `xml:"kml"`
	Document Document `xml:"Document"`
}

type Document struct {
	Folders    []Folder    `xml:"Folder"`
	Placemarks []Placemark `xml:"Placemark"`
}

type Folder struct {
	Name       string      `xml:"name"`
	Placemarks []Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Point       *Point   `xml:"Point"`
	Polygon     *Polygon `xml:"Polygon"`
}

type Point struct {
	Coordinates string `xml:"coordinates"`
}

type Polygon struct {
	OuterBoundaryIs struct {
		LinearRing struct {
			Coordinates string `xml:"coordinates"`
		} `xml:"LinearRing"`
	} `xml:"outerBoundaryIs"`
}

func PopulateLocations(db *sqlx.DB, filePath string) error {
	log.Printf("Opening KML file: %s", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open KML file: %v", err)
	}
	defer file.Close()

	// Parse the KML file
	var kml KML
	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&kml); err != nil {
		return fmt.Errorf("failed to parse KML: %v", err)
	}

	// Start a transaction
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare statement
	stmt, err := tx.Preparex("INSERT INTO locations (name, latitude, longitude) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Process document-level placemarks
	processPlacemarks(kml.Document.Placemarks, stmt)

	// Process placemarks in folders
	for _, folder := range kml.Document.Folders {
		log.Printf("Processing folder: %s", folder.Name)
		processPlacemarks(folder.Placemarks, stmt)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Verify the count
	var count int
	if err := db.Get(&count, "SELECT COUNT(*) FROM locations"); err != nil {
		log.Printf("Error counting locations: %v", err)
	} else {
		log.Printf("Successfully imported %d locations", count)
	}

	return nil
}

func processPlacemarks(placemarks []Placemark, stmt *sqlx.Stmt) {
	for _, p := range placemarks {
		// Skip if name is empty
		if p.Name == "" {
			continue
		}

		// Process Point placemark
		if p.Point != nil && p.Point.Coordinates != "" {
			coords := extractCoordinates(p.Point.Coordinates)
			if coords != nil {
				insertLocation(stmt, p.Name, coords[1], coords[0]) // lat, lon
			}
			continue
		}

		// Process Polygon placemark - use first coordinate as center point
		if p.Polygon != nil {
			coords := p.Polygon.OuterBoundaryIs.LinearRing.Coordinates
			if coords != "" {
				// Take the first coordinate pair as the reference point
				coordPairs := strings.Split(strings.TrimSpace(coords), "\n")
				if len(coordPairs) > 0 {
					firstCoord := extractCoordinates(coordPairs[0])
					if firstCoord != nil {
						insertLocation(stmt, p.Name, firstCoord[1], firstCoord[0]) // lat, lon
					}
				}
			}
		}
	}
}

func extractCoordinates(coordStr string) []float64 {
	parts := strings.Split(strings.TrimSpace(coordStr), ",")
	if len(parts) < 2 {
		return nil
	}

	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		log.Printf("Error parsing longitude: %v", err)
		return nil
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		log.Printf("Error parsing latitude: %v", err)
		return nil
	}

	return []float64{lon, lat}
}

func insertLocation(stmt *sqlx.Stmt, name string, lat, lon float64) {
	_, err := stmt.Exec(name, lat, lon)
	if err != nil {
		log.Printf("Failed to insert location %s: %v", name, err)
	} else {
		log.Printf("Inserted location: Name=%s, Latitude=%f, Longitude=%f", name, lat, lon)
	}
}
