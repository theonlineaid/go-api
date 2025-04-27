package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==========================
// Data Structures
// ==========================

type AttributeInput struct {
	AttributeName string `json:"attribute_name" binding:"required"`
	Status        int    `json:"status" binding:"required"`
}

type AttributeValue struct {
	ID          int    `json:"id"`
	AttributeID int    `json:"attribute_id"`
	Value       string `json:"value"`
	Status      int    `json:"status"`
}

type Attribute struct {
	ID              int              `json:"id"`
	AttributeName   string           `json:"attribute_name"`
	Status          int              `json:"status"`
	AttributeValues []AttributeValue `json:"attribute_values"`
}

// ==========================
// Create Attribute
// ==========================

func CreateAttribute(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AttributeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var attributeID int
		err := db.QueryRow(`
			INSERT INTO attributes (attribute_name, status)
			VALUES ($1, $2)
			RETURNING id
		`, input.AttributeName, input.Status).Scan(&attributeID)

		if err != nil {
			log.Println("Error creating attribute:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Attribute created successfully",
			"attribute_id":   attributeID,
			"attribute_name": input.AttributeName,
		})
	}
}

// ==========================
// Get All Attributes + Values
// ==========================

func GetAllAttributes(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch all attributes
		attrRows, err := db.Query(`SELECT id, attribute_name, status FROM attributes WHERE status = 1`)
		if err != nil {
			log.Println("Error fetching attributes:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer attrRows.Close()

		var attributes []Attribute
		for attrRows.Next() {
			var attr Attribute
			if err := attrRows.Scan(&attr.ID, &attr.AttributeName, &attr.Status); err != nil {
				log.Println("Error scanning attribute:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}

			// Fetch attribute values for this attribute
			valueRows, err := db.Query(`
				SELECT id, attribute_id, value, status
				FROM attribute_values
				WHERE attribute_id = $1 AND status = 1
			`, attr.ID)
			if err != nil {
				log.Println("Error fetching attribute values:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}

			for valueRows.Next() {
				var val AttributeValue
				if err := valueRows.Scan(&val.ID, &val.AttributeID, &val.Value, &val.Status); err != nil {
					log.Println("Error scanning attribute value:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
					return
				}
				attr.AttributeValues = append(attr.AttributeValues, val)
			}
			valueRows.Close()

			attributes = append(attributes, attr)
		}

		if err := attrRows.Err(); err != nil {
			log.Println("Error iterating attributes:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"attributes": attributes,
		})
	}
}

// ==========================
// Get Single Attribute By ID
// ==========================

func GetAttributeByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		attributeID := c.Param("id")

		var attribute Attribute
		err := db.QueryRow(`
			SELECT id, attribute_name, status
			FROM attributes
			WHERE id = $1 AND status = 1
		`, attributeID).Scan(&attribute.ID, &attribute.AttributeName, &attribute.Status)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Attribute not found"})
				return
			}
			log.Println("Error fetching attribute:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		// Fetch attribute values for the attribute
		valueRows, err := db.Query(`
			SELECT id, attribute_id, value, status
			FROM attribute_values
			WHERE attribute_id = $1 AND status = 1
		`, attribute.ID)
		if err != nil {
			log.Println("Error fetching attribute values:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer valueRows.Close()

		for valueRows.Next() {
			var val AttributeValue
			if err := valueRows.Scan(&val.ID, &val.AttributeID, &val.Value, &val.Status); err != nil {
				log.Println("Error scanning attribute value:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}
			attribute.AttributeValues = append(attribute.AttributeValues, val)
		}

		c.JSON(http.StatusOK, attribute)
	}
}

// ==========================
// Update Attribute
// ==========================

func UpdateAttribute(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		attributeID := c.Param("id")

		var input AttributeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec(`
			UPDATE attributes
			SET attribute_name = $1, status = $2
			WHERE id = $3
		`, input.AttributeName, input.Status, attributeID)

		if err != nil {
			log.Println("Error updating attribute:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Attribute updated successfully",
			"attribute_name": input.AttributeName,
			"status":         input.Status,
		})
	}
}

// ==========================
// Delete (Soft Delete) Attribute
// ==========================

func DeleteAttribute(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		attributeID := c.Param("id")

		_, err := db.Exec(`
			UPDATE attributes
			SET status = 0
			WHERE id = $1
		`, attributeID)

		if err != nil {
			log.Println("Error deleting attribute:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Attribute deleted successfully",
		})
	}
}
