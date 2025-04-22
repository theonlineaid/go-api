// handlers/attribute_value.go
package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"my-api/models"
	"net/http"
)

func CreateAttributeValue(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var attrValue models.AttributeValue
		if err := c.BindJSON(&attrValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if attrValue.Value == "" || attrValue.AttributeID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Value and attribute ID are required"})
			return
		}
		if attrValue.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID should not be provided"})
			return
		}

		log.Printf("Received attribute value: %+v", attrValue)

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM attributes WHERE id = $1", attrValue.AttributeID).Scan(&count)
		if err != nil || count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attribute_id"})
			return
		}

		query := `
            INSERT INTO attribute_values (attribute_id, value, status)
            VALUES ($1, $2, $3)
            RETURNING id`
		var attrValueID int
		err = db.QueryRow(query, attrValue.AttributeID, attrValue.Value, attrValue.Status).Scan(&attrValueID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				log.Printf("Duplicate key error: %v", pqErr)
				c.JSON(http.StatusConflict, gin.H{"error": "Attribute value already exists"})
				return
			}
			log.Printf("Insert attribute value error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attribute value"})
			return
		}

		attrValue.ID = attrValueID
		c.JSON(http.StatusCreated, gin.H{"message": "Attribute value created", "attribute_value": attrValue})
	}
}
