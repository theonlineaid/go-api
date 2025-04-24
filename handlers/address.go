package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"my-api/models"

	"github.com/gin-gonic/gin"
)

type AddressInput struct {
	AddressLine1 string `json:"address_line1" binding:"required"`
	City         string `json:"city" binding:"required"`
	Country      string `json:"country" binding:"required"`
	PostalCode   string `json:"postal_code" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=home office other"`
}

type AdminAddressInput struct {
	UserID       int    `json:"user_id" binding:"required,min=1"`
	AddressLine1 string `json:"address_line1" binding:"required"`
	City         string `json:"city" binding:"required"`
	Country      string `json:"country" binding:"required"`
	PostalCode   string `json:"postal_code" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=home office other"`
}

// AddAddress creates a new address for the user
func AddAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var input AddressInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		query := `
			INSERT INTO addresses (user_id, address_line1, city, country, postal_code, type, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			RETURNING id
		`
		var addressID int
		err := db.QueryRow(query, userID, input.AddressLine1, input.City, input.Country, input.PostalCode, input.Type).Scan(&addressID)
		if err != nil {
			log.Println("Error inserting address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Address added",
			"address_id": addressID,
			"type":       input.Type,
		})
	}
}

// GetAddresses retrieves all addresses for the user
func GetAddresses(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		rows, err := db.Query(`
			SELECT id, user_id, address_line1, city, country, postal_code, type, created_at
			FROM addresses
			WHERE user_id = $1
		`, userID)
		if err != nil {
			log.Println("Error querying addresses:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer rows.Close()

		var addresses []models.Address
		for rows.Next() {
			var addr models.Address
			err := rows.Scan(&addr.ID, &addr.UserID, &addr.AddressLine1, &addr.City, &addr.Country,
				&addr.PostalCode, &addr.Type, &addr.CreatedAt)
			if err != nil {
				log.Println("Error scanning address:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}
			addresses = append(addresses, addr)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Addresses retrieved",
			"addresses": addresses,
		})
	}
}

// UpdateAddress updates an existing address
func UpdateAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		addressIDStr := c.Param("id")
		addressID, err := strconv.Atoi(addressIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		var input AddressInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer tx.Rollback()

		var existsAddr bool
		err = tx.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM addresses WHERE id = $1 AND user_id = $2)
		`, addressID, userID).Scan(&existsAddr)
		if err != nil {
			log.Println("Error verifying address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		if !existsAddr {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found or not owned by user"})
			return
		}

		result, err := tx.Exec(`
			UPDATE addresses
			SET address_line1 = $1, city = $2, country = $3, postal_code = $4, type = $5
			WHERE id = $6 AND user_id = $7
		`, input.AddressLine1, input.City, input.Country, input.PostalCode, input.Type, addressID, userID)
		if err != nil {
			log.Println("Error updating address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Address updated",
			"address_id": addressID,
		})
	}
}

func GetAddressByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		addressIDStr := c.Param("id")
		addressID, err := strconv.Atoi(addressIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		var addr models.Address
		err = db.QueryRow(`
			SELECT id, user_id, address_line1, city, country, postal_code, type, created_at
			FROM addresses
			WHERE id = $1 AND user_id = $2
		`, addressID, userID).Scan(&addr.ID, &addr.UserID, &addr.AddressLine1,
			&addr.City, &addr.Country, &addr.PostalCode, &addr.Type, &addr.CreatedAt)
		if err != nil {
			log.Println("Error querying address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Address retrieved",
			"address": addr,
		})
	}
}

// DeleteAddress deletes an address
func DeleteAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		addressIDStr := c.Param("id")
		addressID, err := strconv.Atoi(addressIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		result, err := db.Exec(`
			DELETE FROM addresses
			WHERE id = $1 AND user_id = $2
		`, addressID, userID)
		if err != nil {
			log.Println("Error deleting address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found or not owned by user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Address deleted",
			"address_id": addressID,
		})
	}
}

// AdminCreateAddress creates a new address for any user (admin only)
func AdminCreateAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AdminAddressInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Verify user exists
		var existsUser bool
		err := db.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)
		`, input.UserID).Scan(&existsUser)
		if err != nil {
			log.Println("Error verifying user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		if !existsUser {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		query := `
			INSERT INTO addresses (user_id, address_line1, city, country, postal_code, type, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			RETURNING id
		`
		var addressID int
		err = db.QueryRow(query, input.UserID, input.AddressLine1, input.City, input.Country, input.PostalCode, input.Type).Scan(&addressID)
		if err != nil {
			log.Println("Error inserting address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Address created",
			"address_id": addressID,
		})
	}
}

// AdminUpdateAddress updates an existing address by ID (admin only)
func AdminUpdateAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		addressIDStr := c.Param("id")
		addressID, err := strconv.Atoi(addressIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		var input AddressInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer tx.Rollback()

		// Verify address exists
		var existsAddr bool
		err = tx.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM addresses WHERE id = $1)
		`, addressID).Scan(&existsAddr)
		if err != nil {
			log.Println("Error verifying address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		if !existsAddr {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}

		result, err := tx.Exec(`
			UPDATE addresses
			SET address_line1 = $1, city = $2, country = $3, postal_code = $4, type = $5
			WHERE id = $6
		`, input.AddressLine1, input.City, input.Country, input.PostalCode, input.Type, addressID)
		if err != nil {
			log.Println("Error updating address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Address updated",
			"address_id": addressID,
		})
	}
}

// Existing handlers (unchanged)
func AddShippingAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var input AddressInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer tx.Rollback()

		var defaultExists bool
		err = tx.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM shipping_addresses WHERE user_id = $1 AND is_default = true)
		`, userID).Scan(&defaultExists)
		if err != nil {
			log.Println("Error checking default shipping address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		query := `
			INSERT INTO shipping_addresses (user_id, address_line1, city, country, postal_code, type, is_default, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			RETURNING id
		`
		isDefault := !defaultExists
		var addressID int
		err = tx.QueryRow(query, userID, input.AddressLine1, input.City, input.Country, input.PostalCode, input.Type, isDefault).Scan(&addressID)
		if err != nil {
			log.Println("Error inserting shipping address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Shipping address added",
			"address_id": addressID,
			"is_default": isDefault,
			"type":       input.Type,
		})
	}
}

func AddBillingAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var input AddressInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer tx.Rollback()

		var defaultExists bool
		err = tx.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM billing_addresses WHERE user_id = $1 AND is_default = true)
		`, userID).Scan(&defaultExists)
		if err != nil {
			log.Println("Error checking default billing address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		query := `
			INSERT INTO billing_addresses (user_id, address_line1, city, country, postal_code, type, is_default, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			RETURNING id
		`
		isDefault := !defaultExists
		var addressID int
		err = tx.QueryRow(query, userID, input.AddressLine1, input.City, input.Country, input.PostalCode, input.Type, isDefault).Scan(&addressID)
		if err != nil {
			log.Println("Error inserting billing address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Billing address added",
			"address_id": addressID,
			"is_default": isDefault,
			"type":       input.Type,
		})
	}
}

func SetDefaultShippingAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		addressIDStr := c.Param("id")
		addressID, err := strconv.Atoi(addressIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer tx.Rollback()

		var existsAddr bool
		err = tx.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM shipping_addresses WHERE id = $1 AND user_id = $2)
		`, addressID, userID).Scan(&existsAddr)
		if err != nil {
			log.Println("Error verifying shipping address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		if !existsAddr {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found or not owned by user"})
			return
		}

		_, err = tx.Exec(`
			UPDATE shipping_addresses SET is_default = FALSE
			WHERE user_id = $1 AND is_default = TRUE
		`, userID)
		if err != nil {
			log.Println("Error unsetting default shipping address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		result, err := tx.Exec(`
			UPDATE shipping_addresses SET is_default = TRUE
			WHERE id = $1 AND user_id = $2
		`, addressID, userID)
		if err != nil {
			log.Println("Error setting default shipping address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Shipping address set as default",
			"address_id": addressID,
		})
	}
}

func SetDefaultBillingAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		addressIDStr := c.Param("id")
		addressID, err := strconv.Atoi(addressIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer tx.Rollback()

		var existsAddr bool
		err = tx.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM billing_addresses WHERE id = $1 AND user_id = $2)
		`, addressID, userID).Scan(&existsAddr)
		if err != nil {
			log.Println("Error verifying billing address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		if !existsAddr {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found or not owned by user"})
			return
		}

		_, err = tx.Exec(`
			UPDATE billing_addresses SET is_default = FALSE
			WHERE user_id = $1 AND is_default = TRUE
		`, userID)
		if err != nil {
			log.Println("Error unsetting default billing address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		result, err := tx.Exec(`
			UPDATE billing_addresses SET is_default = TRUE
			WHERE id = $1 AND user_id = $2
		`, addressID, userID)
		if err != nil {
			log.Println("Error setting default billing address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Billing address set as default",
			"address_id": addressID,
		})
	}
}
