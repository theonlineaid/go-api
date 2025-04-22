// handlers/product.go
package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"my-api/models"
	"net/http"
)

type CreateProductRequest struct {
	Product    models.Product            `json:"product"`
	Attributes []models.ProductAttribute `json:"attributes"`
	Variations []models.VariationProduct `json:"variations"`
}

func AddProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateProductRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		product := req.Product
		// Validate required fields
		if product.ProductName == "" || product.ProductCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product name and code are required"})
			return
		}
		if product.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID should not be provided"})
			return
		}

		// Validate foreign keys
		var count int
		if product.BrandID != 0 {
			err := db.QueryRow("SELECT COUNT(*) FROM brands WHERE id = $1", product.BrandID).Scan(&count)
			if err != nil || count == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand_id"})
				return
			}
		}
		err := db.QueryRow("SELECT COUNT(*) FROM categories WHERE id = $1", product.CategoryID).Scan(&count)
		if err != nil || count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
			return
		}
		if product.SubCategoryID != nil {
			err = db.QueryRow("SELECT COUNT(*) FROM subcategories WHERE id = $1", *product.SubCategoryID).Scan(&count)
			if err != nil || count == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sub_category_id"})
				return
			}
		}

		// Log the received product
		log.Printf("Received product: %+v, attributes: %+v, variations: %+v", product, req.Attributes, req.Variations)

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Transaction begin error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer tx.Rollback()

		// Insert product
		query := `
            INSERT INTO products (
                brand_id, category_id, sub_category_id, p_code, weight, product_name, 
                product_code, price, tax_type, image, serial_number, is_saleable, 
                is_barcode, is_warranty, is_variation, status, rating, current_stock
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
            RETURNING id`
		var productID int
		err = tx.QueryRow(
			query,
			product.BrandID,
			product.CategoryID,
			product.SubCategoryID,
			product.PCode,
			product.Weight,
			product.ProductName,
			product.ProductCode,
			product.Price,
			product.TaxType,
			product.Image,
			product.SerialNumber,
			product.IsSaleable,
			product.IsBarcode,
			product.IsWarranty,
			product.IsVariation,
			product.Status,
			product.Rating,
			product.CurrentStock,
		).Scan(&productID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				log.Printf("Duplicate key error: %v", pqErr)
				c.JSON(http.StatusConflict, gin.H{"error": "Product already exists"})
				return
			}
			log.Printf("Insert product error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		// Insert attributes
		for _, attr := range req.Attributes {
			// Validate attribute_id and attribute_value_id
			err = tx.QueryRow("SELECT COUNT(*) FROM attributes WHERE id = $1", attr.AttributeID).Scan(&count)
			if err != nil || count == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attribute_id"})
				return
			}
			err = tx.QueryRow("SELECT COUNT(*) FROM attribute_values WHERE id = $1 AND attribute_id = $2", attr.AttributeValueID, attr.AttributeID).Scan(&count)
			if err != nil || count == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attribute_value_id"})
				return
			}

			attrQuery := `
                INSERT INTO product_attributes (product_id, attribute_id, attribute_value_id)
                VALUES ($1, $2, $3)`
			_, err = tx.Exec(attrQuery, productID, attr.AttributeID, attr.AttributeValueID)
			if err != nil {
				log.Printf("Insert attribute error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attributes"})
				return
			}
		}

		// Insert variations
		for _, variation := range req.Variations {
			if variation.SKU == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Variation SKU is required"})
				return
			}
			varQuery := `
                INSERT INTO variation_products (product_id, sku, sale_price, default_sell_price, discount, image, current_stock)
                VALUES ($1, $2, $3, $4, $5, $6, $7)`
			_, err = tx.Exec(
				varQuery,
				productID,
				variation.SKU,
				variation.SalePrice,
				variation.DefaultSellPrice,
				variation.Discount,
				variation.Image,
				variation.CurrentStock,
			)
			if err != nil {
				log.Printf("Insert variation error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add variations"})
				return
			}
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			log.Printf("Transaction commit error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		product.ID = productID
		c.JSON(http.StatusCreated, gin.H{
			"message":    "Product added successfully",
			"product":    product,
			"attributes": req.Attributes,
			"variations": req.Variations,
		})
	}
}

func GetProductByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Println("Fetching product with ID:", id)

		var product models.Product
		query := `
            SELECT id, brand_id, category_id, sub_category_id, p_code, weight, product_name, 
                   product_code, price, tax_type, image, serial_number, is_saleable, 
                   is_barcode, is_warranty, is_variation, status, rating, current_stock
            FROM products WHERE id = $1`
		err := db.QueryRow(query, id).Scan(
			&product.ID,
			&product.BrandID,
			&product.CategoryID,
			&product.SubCategoryID,
			&product.PCode,
			&product.Weight,
			&product.ProductName,
			&product.ProductCode,
			&product.Price,
			&product.TaxType,
			&product.Image,
			&product.SerialNumber,
			&product.IsSaleable,
			&product.IsBarcode,
			&product.IsWarranty,
			&product.IsVariation,
			&product.Status,
			&product.Rating,
			&product.CurrentStock,
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

		// Fetch attributes
		var attributes []models.ProductAttribute
		attrQuery := `
            SELECT id, product_id, attribute_id, attribute_value_id
            FROM product_attributes WHERE product_id = $1`
		rows, err := db.Query(attrQuery, id)
		if err != nil {
			log.Println("Get attributes error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attributes"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var attr models.ProductAttribute
			if err := rows.Scan(&attr.ID, &attr.ProductID, &attr.AttributeID, &attr.AttributeValueID); err != nil {
				log.Println("Scan attributes error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan attributes"})
				return
			}
			attributes = append(attributes, attr)
		}

		// Fetch variations
		var variations []models.VariationProduct
		varQuery := `
            SELECT id, product_id, sku, sale_price, default_sell_price, discount, image, current_stock
            FROM variation_products WHERE product_id = $1`
		rows, err = db.Query(varQuery, id)
		if err != nil {
			log.Println("Get variations error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get variations"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var variation models.VariationProduct
			if err := rows.Scan(
				&variation.ID,
				&variation.ProductID,
				&variation.SKU,
				&variation.SalePrice,
				&variation.DefaultSellPrice,
				&variation.Discount,
				&variation.Image,
				&variation.CurrentStock,
			); err != nil {
				log.Println("Scan variations error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan variations"})
				return
			}
			variations = append(variations, variation)
		}

		c.JSON(http.StatusOK, gin.H{
			"product":    product,
			"attributes": attributes,
			"variations": variations,
		})
	}
}
