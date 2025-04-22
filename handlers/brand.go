// handlers/brand.go
package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"my-api/models"
	"net/http"
)

func CreateBrand(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var brand models.Brand
		if err := c.BindJSON(&brand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if brand.BrandName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Brand name is required"})
			return
		}
		if brand.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID should not be provided"})
			return
		}

		log.Printf("Received brand: %+v", brand)

		query := `
            INSERT INTO brands (brand_name, image, status)
            VALUES ($1, $2, $3)
            RETURNING id`
		var brandID int
		err := db.QueryRow(query, brand.BrandName, brand.Image, brand.Status).Scan(&brandID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				log.Printf("Duplicate key error: %v", pqErr)
				c.JSON(http.StatusConflict, gin.H{"error": "Brand already exists"})
				return
			}
			log.Printf("Insert brand error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create brand"})
			return
		}

		brand.ID = brandID
		c.JSON(http.StatusCreated, gin.H{"message": "Brand created", "brand": brand})
	}
}
