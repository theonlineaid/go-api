// handlers/attribute.go
package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"my-api/models"
	"net/http"
)

func CreateAttribute(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var attribute models.Attribute
		if err := c.BindJSON(&attribute); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if attribute.AttributeName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Attribute name is required"})
			return
		}
		if attribute.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID should not be provided"})
			return
		}

		log.Printf("Received attribute: %+v", attribute)

		query := `
            INSERT INTO attributes (attribute_name, status)
            VALUES ($1, $2)
            RETURNING id`
		var attributeID int
		err := db.QueryRow(query, attribute.AttributeName, attribute.Status).Scan(&attributeID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				log.Printf("Duplicate key error: %v", pqErr)
				c.JSON(http.StatusConflict, gin.H{"error": "Attribute already exists"})
				return
			}
			log.Printf("Insert attribute error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attribute"})
			return
		}

		attribute.ID = attributeID
		c.JSON(http.StatusCreated, gin.H{"message": "Attribute created", "attribute": attribute})
	}
}
