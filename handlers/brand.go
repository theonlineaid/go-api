package handlers

import (
	"database/sql"
	"log"
	"my-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBrand(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var brand models.Brand
		if err := c.BindJSON(&brand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand JSON"})
			return
		}

		query := `
        INSERT INTO brands (
            brand_name, image, status,
            is_feature, is_publish, is_special,
            is_approved_by_admin, is_visible_to_guest,
            created_by, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id;
        `

		var brandID int
		err := db.QueryRow(query,
			brand.BrandName, brand.Image, brand.Status,
			brand.IsFeature, brand.IsPublish, brand.IsSpecial,
			brand.IsApprovedByAdmin, brand.IsVisibleToGuest,
			brand.CreatedBy,
		).Scan(&brandID)

		if err != nil {
			log.Println("Insert brand error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert brand"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Brand added successfully",
			"id":      brandID,
		})
	}
}

func GetAllBrands(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
            SELECT id, brand_name, image, status,
                   is_feature, is_publish, is_special,
                   is_approved_by_admin, is_visible_to_guest,
                   created_by, created_at, updated_at
            FROM brands`
		rows, err := db.Query(query)
		if err != nil {
			log.Println("Error querying brands:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch brands"})
			return
		}
		defer rows.Close()

		var brands []models.Brand
		for rows.Next() {
			var brand models.Brand
			err := rows.Scan(
				&brand.ID, &brand.BrandName, &brand.Image, &brand.Status,
				&brand.IsFeature, &brand.IsPublish, &brand.IsSpecial,
				&brand.IsApprovedByAdmin, &brand.IsVisibleToGuest,
				&brand.CreatedBy, &brand.CreatedAt, &brand.UpdatedAt,
			)
			if err != nil {
				log.Println("Error scanning brand:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process brands"})
				return
			}
			brands = append(brands, brand)
		}

		if err = rows.Err(); err != nil {
			log.Println("Error iterating brands:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch brands"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"brands": brands})
	}
}

func GetBrandByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var brand models.Brand
		query := `
            SELECT id, brand_name, image, status,
                   is_feature, is_publish, is_special,
                   is_approved_by_admin, is_visible_to_guest,
                   created_by, created_at, updated_at
            FROM brands WHERE id = $1`
		err := db.QueryRow(query, id).Scan(
			&brand.ID, &brand.BrandName, &brand.Image, &brand.Status,
			&brand.IsFeature, &brand.IsPublish, &brand.IsSpecial,
			&brand.IsApprovedByAdmin, &brand.IsVisibleToGuest,
			&brand.CreatedBy, &brand.CreatedAt, &brand.UpdatedAt,
		)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}
		if err != nil {
			log.Println("Error querying brand:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch brand"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"brand": brand})
	}
}

func UpdateBrand(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var brand models.Brand
		if err := c.BindJSON(&brand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand JSON"})
			return
		}

		query := `
            UPDATE brands SET
                brand_name = $1,
                image = $2,
                status = $3,
                is_feature = $4,
                is_publish = $5,
                is_special = $6,
                is_approved_by_admin = $7,
                is_visible_to_guest = $8,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = $9
            RETURNING id`
		var updatedID int
		err := db.QueryRow(query,
			brand.BrandName, brand.Image, brand.Status,
			brand.IsFeature, brand.IsPublish, brand.IsSpecial,
			brand.IsApprovedByAdmin, brand.IsVisibleToGuest,
			id,
		).Scan(&updatedID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}
		if err != nil {
			log.Println("Error updating brand:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update brand"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Brand updated successfully",
			"id":      updatedID,
		})
	}
}

func DeleteBrand(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		query := `DELETE FROM brands WHERE id = $1 RETURNING id`
		var deletedID int
		err := db.QueryRow(query, id).Scan(&deletedID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}
		if err != nil {
			log.Println("Error deleting brand:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete brand"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Brand deleted successfully",
			"id":      deletedID,
		})
	}
}
