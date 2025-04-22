// handlers/category.go
package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"my-api/models"
	"net/http"
)

func CreateCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		if err := c.BindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if category.CategoryName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
			return
		}
		if category.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID should not be provided"})
			return
		}

		log.Printf("Received category: %+v", category)

		query := `
            INSERT INTO categories (category_name, image, is_special, price_visibility, status)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id`
		var categoryID int
		err := db.QueryRow(
			query,
			category.CategoryName,
			category.Image,
			category.IsSpecial,
			category.PriceVisibility,
			category.Status,
		).Scan(&categoryID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				log.Printf("Duplicate key error: %v", pqErr)
				c.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
				return
			}
			log.Printf("Insert category error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			return
		}

		category.ID = categoryID
		c.JSON(http.StatusCreated, gin.H{"message": "Category created", "category": category})
	}
}
