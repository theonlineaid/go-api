package handlers

import (
	"database/sql"
	"log"
	"my-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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
                is_special, is_featured, is_approved, is_published, position,
                price_visibility, status, created_by, created_at, updated_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
            RETURNING id
        `

		var categoryID int
		err := db.QueryRow(query,
			category.Code, category.CategoryName, category.CategoryImg, category.Image,
			category.CategoryVisibility, category.IsSpecial, category.IsFeatured,
			category.IsApproved, category.IsPublished, category.Position,
			category.PriceVisibility, category.Status, category.CreatedBy, now, now,
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
		category.SubCategories = []models.SubCategory{}

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
                code = $1, category_name = $2, category_img = $3, image = $4,
                category_visibility = $5, is_special = $6, is_featured = $7,
                is_approved = $8, is_published = $9, position = $10,
                price_visibility = $11, status = $12, updated_at = $13
            WHERE id = $14
            RETURNING id
        `

		var updatedID int
		err := db.QueryRow(query,
			category.Code, category.CategoryName, category.CategoryImg, category.Image,
			category.CategoryVisibility, category.IsSpecial, category.IsFeatured,
			category.IsApproved, category.IsPublished, category.Position,
			category.PriceVisibility, category.Status, now, id,
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
                   is_special, is_featured, is_approved, is_published, position,
                   price_visibility, status, created_by, created_at, updated_at
            FROM categories WHERE id = $1`
		err := db.QueryRow(query, id).Scan(
			&category.ID, &category.Code, &category.CategoryName, &category.CategoryImg,
			&category.Image, &category.CategoryVisibility, &category.IsSpecial,
			&category.IsFeatured, &category.IsApproved, &category.IsPublished,
			&category.Position, &category.PriceVisibility, &category.Status,
			&category.CreatedBy, &category.CreatedAt, &category.UpdatedAt,
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

		// Fetch subcategories
		subQuery := `
            SELECT id, category_id, subcategory_name, image, status
            FROM subcategories WHERE category_id = $1`
		rows, err := db.Query(subQuery, id)
		if err != nil {
			log.Printf("Error querying subcategories: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcategories"})
			return
		}
		defer rows.Close()

		category.SubCategories = []models.SubCategory{}
		for rows.Next() {
			var subCategory models.SubCategory
			err := rows.Scan(
				&subCategory.ID, &subCategory.CategoryID, &subCategory.SubCategoryName,
				&subCategory.Image, &subCategory.Status,
			)
			if err != nil {
				log.Printf("Error scanning subcategory: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process subcategories"})
				return
			}
			category.SubCategories = append(category.SubCategories, subCategory)
		}

		category.ProductsCount = 0 // Update with actual count if needed

		c.JSON(http.StatusOK, gin.H{"category": category})
	}
}

func GetAllCategories(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
            SELECT id, code, category_name, category_img, image, category_visibility,
                   is_special, is_featured, is_approved, is_published, position,
                   price_visibility, status, created_by, created_at, updated_at
            FROM categories`
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Error querying categories: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}
		defer rows.Close()

		var categories []models.Category
		categoryIDs := []int{}
		for rows.Next() {
			var category models.Category
			err := rows.Scan(
				&category.ID, &category.Code, &category.CategoryName, &category.CategoryImg,
				&category.Image, &category.CategoryVisibility, &category.IsSpecial,
				&category.IsFeatured, &category.IsApproved, &category.IsPublished,
				&category.Position, &category.PriceVisibility, &category.Status,
				&category.CreatedBy, &category.CreatedAt, &category.UpdatedAt,
			)
			if err != nil {
				log.Printf("Error scanning category: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process categories"})
				return
			}
			category.SubCategories = []models.SubCategory{}
			category.ProductsCount = 0
			categories = append(categories, category)
			categoryIDs = append(categoryIDs, category.ID)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error iterating categories: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

		// Fetch subcategories for all categories
		if len(categoryIDs) > 0 {
			subQuery := `
                SELECT id, category_id, subcategory_name, image, status
                FROM subcategories WHERE category_id = ANY($1)`
			rows, err := db.Query(subQuery, pq.Array(categoryIDs))
			if err != nil {
				log.Printf("Error querying subcategories: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcategories"})
				return
			}
			defer rows.Close()

			subCategoriesMap := make(map[int][]models.SubCategory)
			for rows.Next() {
				var subCategory models.SubCategory
				err := rows.Scan(
					&subCategory.ID, &subCategory.CategoryID, &subCategory.SubCategoryName,
					&subCategory.Image, &subCategory.Status,
				)
				if err != nil {
					log.Printf("Error scanning subcategory: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process subcategories"})
					return
				}
				subCategoriesMap[subCategory.CategoryID] = append(subCategoriesMap[subCategory.CategoryID], subCategory)
			}

			for i := range categories {
				categories[i].SubCategories = subCategoriesMap[categories[i].ID]
			}
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}
