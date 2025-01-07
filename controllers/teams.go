package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Team struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Leader string `json:"leader"`
}

func GetTeams(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var teams []Team
		err := db.Select(&teams, "SELECT * FROM teams")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
			return
		}
		c.JSON(http.StatusOK, teams)
	}
}

func AddTeam(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var team Team
		if err := c.ShouldBindJSON(&team); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec("INSERT INTO teams (name, leader) VALUES (?, ?)", team.Name, team.Leader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add team"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Team added successfully"})
	}
}

// Delete a team
func DeleteTeam(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM teams WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
	}
}

// Update a team
func UpdateTeam(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var team Team
		if err := c.ShouldBindJSON(&team); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec("UPDATE teams SET name = ?, leader = ? WHERE id = ?", team.Name, team.Leader, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Team updated successfully"})
	}
}
