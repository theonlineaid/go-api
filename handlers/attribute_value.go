// handlers/attribute_value.go
package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// AttributeValueInput for creating or updating attribute values
type AttributeValueInput struct {
	AttributeID int    `json:"attribute_id" binding:"required"`
	Value       string `json:"value" binding:"required"`
	Status      int    `json:"status" binding:"required"`
}

// AttributeValueResponse represents a single attribute value
type AttributeValueResponse struct {
	ID          int    `json:"id"`
	AttributeID int    `json:"attribute_id"`
	Value       string `json:"value"`
	Status      int    `json:"status"`
}

// CreateAttributeValue creates a new attribute value
// CreateAttributeValue creates a new attribute value
func CreateAttributeValue(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AttributeValueInput

		// Bind JSON input
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Step 1: Check if attribute_id exists
		var exists bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM attributes WHERE id = $1)`, input.AttributeID).Scan(&exists)
		if err != nil {
			log.Println("Error checking attribute_id existence:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attribute_id, attribute does not exist"})
			return
		}

		// Step 2: Insert attribute value
		var valueID int
		err = db.QueryRow(`
			INSERT INTO attribute_values (attribute_id, value, status)
			VALUES ($1, $2, $3)
			RETURNING id
		`, input.AttributeID, input.Value, input.Status).Scan(&valueID)

		if err != nil {
			// Check if foreign key error
			if pgErr, ok := err.(*pq.Error); ok {
				if pgErr.Code == "23503" { // foreign_key_violation
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attribute_id"})
					return
				}
			}
			log.Println("Error creating attribute value:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Attribute value created successfully",
			"value_id":     valueID,
			"attribute_id": input.AttributeID,
			"value":        input.Value,
		})
	}
}

// GetAllAttributeValues retrieves all attribute values for a given attribute ID
func GetAllAttributeValues(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		attributeID := c.Param("attribute_id")

		rows, err := db.Query(`
			SELECT id, attribute_id, value, status
			FROM attribute_values
			WHERE attribute_id = $1
		`, attributeID)
		if err != nil {
			log.Println("Error fetching attribute values:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attribute values"})
			return
		}
		defer rows.Close()

		var attributeValues []AttributeValueResponse

		for rows.Next() {
			var av AttributeValueResponse
			if err := rows.Scan(&av.ID, &av.AttributeID, &av.Value, &av.Status); err != nil {
				log.Println("Error scanning attribute value:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan attribute values"})
				return
			}
			attributeValues = append(attributeValues, av)
		}

		c.JSON(http.StatusOK, gin.H{
			"attribute_values": attributeValues,
		})
	}
}

// GetAttributeValue retrieves a single attribute value by ID
func GetAttributeValue(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var av AttributeValueResponse
		err := db.QueryRow(`
			SELECT id, attribute_id, value, status
			FROM attribute_values
			WHERE id = $1
		`, id).Scan(&av.ID, &av.AttributeID, &av.Value, &av.Status)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Attribute value not found"})
				return
			}
			log.Println("Error fetching attribute value:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attribute value"})
			return
		}

		c.JSON(http.StatusOK, av)
	}
}

// UpdateAttributeValue updates an existing attribute value
func UpdateAttributeValue(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input struct {
			Value  string `json:"value" binding:"required"`
			Status int    `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		result, err := db.Exec(`
			UPDATE attribute_values
			SET value = $1, status = $2
			WHERE id = $3
		`, input.Value, input.Status, id)
		if err != nil {
			log.Println("Error updating attribute value:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attribute value"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Attribute value not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Attribute value updated successfully",
		})
	}
}

// DeleteAttributeValue deletes an attribute value by ID
func DeleteAttributeValue(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := db.Exec(`
			DELETE FROM attribute_values WHERE id = $1
		`, id)
		if err != nil {
			log.Println("Error deleting attribute value:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attribute value"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Attribute value not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Attribute value deleted successfully",
		})
	}
}
