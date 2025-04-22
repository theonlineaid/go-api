// handlers/subcategory.go
package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"my-api/models"
	"net/http"
)

func CreateSubcategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var subcategory models.Subcategory
		if err := c.BindJSON(&subcategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if subcategory.SubcategoryName == "" || subcategory.CategoryID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subcategory name and category ID are required"})
			return
		}
		if subcategory.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID should not be provided"})
			return
		}

		log.Printf("Received subcategory: %+v", subcategory)

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM categories WHERE id = $1", subcategory.CategoryID).Scan(&count)
		if err != nil || count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
			return
		}

		query := `
            INSERT INTO subcategories (category_id, subcategory_name, image, status)
            VALUES ($1, $2, $3, $4)
            RETURNING id`
		var subcategoryID int
		err = db.QueryRow(
			query,
			subcategory.CategoryID,
			subcategory.SubcategoryName,
			subcategory.Image,
			subcategory.Status,
		).Scan(&subcategoryID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				log.Printf("Duplicate key error: %v", pqErr)
				c.JSON(http.StatusConflict, gin.H{"error": "Subcategory already exists"})
				return
			}
			log.Printf("Insert subcategory error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subcategory"})
			return
		}

		subcategory.ID = subcategoryID
		c.JSON(http.StatusCreated, gin.H{"message": "Subcategory created", "subcategory": subcategory})
	}
}
