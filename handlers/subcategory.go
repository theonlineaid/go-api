package handlers

import (
	"database/sql"
	"log"
	"my-api/models"
	"net/http"
	_ "time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func CreateSubCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var subCategory models.SubCategory
		if err := c.ShouldBindJSON(&subCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
			return
		}

		if subCategory.SubCategoryName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subcategory name is required"})
			return
		}
		if subCategory.CategoryID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID is required"})
			return
		}

		query := `
            INSERT INTO subcategories (
                category_id, subcategory_name, image, status
            ) VALUES ($1, $2, $3, $4)
            RETURNING id
        `

		var subCategoryID int
		err := db.QueryRow(query,
			subCategory.CategoryID, subCategory.SubCategoryName, subCategory.Image, subCategory.Status,
		).Scan(&subCategoryID)
		if err != nil {
			log.Printf("Failed to insert subcategory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subcategory"})
			return
		}

		subCategory.ID = subCategoryID
		c.JSON(http.StatusCreated, gin.H{
			"message":     "Subcategory created successfully",
			"subcategory": subCategory,
		})
	}
}

func UpdateSubCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var subCategory models.SubCategory
		if err := c.ShouldBindJSON(&subCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
			return
		}

		if subCategory.SubCategoryName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subcategory name is required"})
			return
		}

		query := `
            UPDATE subcategories SET
                subcategory_name = $1,
                image = $2,
                status = $3
            WHERE id = $4
            RETURNING id, category_id
        `

		var updatedID, categoryID int
		err := db.QueryRow(query,
			subCategory.SubCategoryName, subCategory.Image, subCategory.Status, id,
		).Scan(&updatedID, &categoryID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found"})
			return
		}
		if err != nil {
			log.Printf("Error updating subcategory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subcategory"})
			return
		}

		subCategory.ID = updatedID
		subCategory.CategoryID = categoryID
		c.JSON(http.StatusOK, gin.H{
			"message":     "Subcategory updated successfully",
			"subcategory": subCategory,
		})
	}
}

func DeleteSubCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		query := `DELETE FROM subcategories WHERE id = $1 RETURNING id`
		var deletedID int
		err := db.QueryRow(query, id).Scan(&deletedID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found"})
			return
		}
		if err != nil {
			log.Printf("Error deleting subcategory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subcategory"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Subcategory deleted successfully",
			"id":      deletedID,
		})
	}
}

func GetAllSubCategories(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
            SELECT id, category_id, subcategory_name, image, status
            FROM subcategories`
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Error querying subcategories: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcategories"})
			return
		}
		defer rows.Close()

		var subCategories []models.SubCategory
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
			subCategories = append(subCategories, subCategory)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error iterating subcategories: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcategories"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"subcategories": subCategories})
	}
}

func GetSubCategoryByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var subCategory models.SubCategory
		query := `
            SELECT id, category_id, subcategory_name, image, status
            FROM subcategories WHERE id = $1`
		err := db.QueryRow(query, id).Scan(
			&subCategory.ID, &subCategory.CategoryID, &subCategory.SubCategoryName,
			&subCategory.Image, &subCategory.Status,
		)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found"})
			return
		}
		if err != nil {
			log.Printf("Error querying subcategory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcategory"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"subcategory": subCategory})
	}
}
