package handlers

import (
	"database/sql"
	"log"
	"my-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
			return
		}

		if category.CategoryName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
			return
		}
		if category.CreatedBy == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CreatedBy is required"})
			return
		}

		now := time.Now()
		query := `
            INSERT INTO categories (
                code, category_name, category_img, image, category_visibility,
                is_special, is_featured, position, price_visibility, status,
                created_by, created_at, updated_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
            RETURNING id
        `

		var categoryID int
		err := db.QueryRow(query,
			category.Code, category.CategoryName, category.CategoryImg, category.Image,
			category.CategoryVisibility, category.IsSpecial, category.IsFeatured,
			category.Position, category.PriceVisibility, category.Status,
			category.CreatedBy, now, now,
		).Scan(&categoryID)
		if err != nil {
			log.Printf("Failed to insert category: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			return
		}

		category.ID = categoryID
		category.CreatedAt = now
		category.UpdatedAt = now
		category.ProductsCount = 0
		category.SubCategories = []any{}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "Category created successfully",
			"category": category,
		})
	}
}

func UpdateCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
			return
		}

		if category.CategoryName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
			return
		}

		now := time.Now()
		query := `
            UPDATE categories SET
                code = $1,
                category_name = $2,
                category_img = $3,
                image = $4,
                category_visibility = $5,
                is_special = $6,
                is_featured = $7,
                position = $8,
                price_visibility = $9,
                status = $10,
                updated_at = $11
            WHERE id = $12
            RETURNING id
        `

		var updatedID int
		err := db.QueryRow(query,
			category.Code, category.CategoryName, category.CategoryImg, category.Image,
			category.CategoryVisibility, category.IsSpecial, category.IsFeatured,
			category.Position, category.PriceVisibility, category.Status,
			now, id,
		).Scan(&updatedID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		if err != nil {
			log.Printf("Error updating category: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Category updated successfully",
			"id":      updatedID,
		})
	}
}

func DeleteCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		query := `DELETE FROM categories WHERE id = $1 RETURNING id`
		var deletedID int
		err := db.QueryRow(query, id).Scan(&deletedID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		if err != nil {
			log.Printf("Error deleting category: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Category deleted successfully",
			"id":      deletedID,
		})
	}
}

func GetCategoryByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category
		query := `
            SELECT id, code, category_name, category_img, image, category_visibility,
                   is_special, is_featured, position, price_visibility, status,
                   created_by, created_at, updated_at
            FROM categories WHERE id = $1`
		err := db.QueryRow(query, id).Scan(
			&category.ID, &category.Code, &category.CategoryName, &category.CategoryImg,
			&category.Image, &category.CategoryVisibility, &category.IsSpecial,
			&category.IsFeatured, &category.Position, &category.PriceVisibility,
			&category.Status, &category.CreatedBy, &category.CreatedAt, &category.UpdatedAt,
		)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		if err != nil {
			log.Printf("Error querying category: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
			return
		}

		category.ProductsCount = 0
		category.SubCategories = []any{}

		c.JSON(http.StatusOK, gin.H{"category": category})
	}
}

func GetAllCategories(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
            SELECT id, code, category_name, category_img, image, category_visibility,
                   is_special, is_featured, position, price_visibility, status,
                   created_by, created_at, updated_at
            FROM categories`
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Error querying categories: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}
		defer rows.Close()

		var categories []models.Category
		for rows.Next() {
			var category models.Category
			err := rows.Scan(
				&category.ID, &category.Code, &category.CategoryName, &category.CategoryImg,
				&category.Image, &category.CategoryVisibility, &category.IsSpecial,
				&category.IsFeatured, &category.Position, &category.PriceVisibility,
				&category.Status, &category.CreatedBy, &category.CreatedAt, &category.UpdatedAt,
			)
			if err != nil {
				log.Printf("Error scanning category: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process categories"})
				return
			}
			category.ProductsCount = 0
			category.SubCategories = []any{}
			categories = append(categories, category)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error iterating categories: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}
