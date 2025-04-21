package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"my-api/models"

	"github.com/gin-gonic/gin"
)

func AddProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product JSON"})
			return
		}

		// Check if the category exists
		var categoryExists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", product.CategoryID).Scan(&categoryExists)
		if err != nil {
			log.Println("Error checking category:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check category"})
			return
		}

		if !categoryExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category does not exist"})
			return
		}

		// Insert the product
		query := `INSERT INTO products (id, brand_id, category_id, product_name, product_code, tax_type, image, serial_number, is_saleable, is_barcode, is_warranty, is_variation, status, rating, current_stock)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

		_, err = db.Exec(query,
			product.ID, product.BrandID, product.CategoryID, product.ProductName, product.ProductCode, product.TaxType, product.Image,
			product.SerialNumber, product.IsSaleable, product.IsBarcode, product.IsWarranty, product.IsVariation, product.Status,
			product.Rating, product.CurrentStock)

		if err != nil {
			log.Println("Insert product error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully"})
	}
}

// handlers/product.go
func GetProductByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		log.Println("Fetching product with ID:", id)

		var product models.Product
		query := `SELECT id, brand_id, category_id, product_name, product_code, tax_type, image,
				  serial_number, is_saleable, is_barcode, is_warranty, is_variation, status, rating, current_stock
				  FROM products WHERE id = $1`

		err := db.QueryRow(query, id).Scan(
			&product.ID, &product.BrandID, &product.CategoryID, &product.ProductName, &product.ProductCode,
			&product.TaxType, &product.Image, &product.SerialNumber, &product.IsSaleable, &product.IsBarcode,
			&product.IsWarranty, &product.IsVariation, &product.Status, &product.Rating, &product.CurrentStock,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			} else {
				log.Println("Get product error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product"})
			}
			return
		}

		c.JSON(http.StatusOK, product)
	}
}
