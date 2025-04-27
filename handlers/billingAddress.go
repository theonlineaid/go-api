package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func GetBillingAddresses(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		rows, err := db.Query(`
			SELECT id, address_line1, city, country, postal_code, type, is_default
			FROM billing_addresses
			WHERE user_id = $1
		`, userID)
		if err != nil {
			log.Println("Error fetching billing addresses:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer rows.Close()

		var addresses []gin.H
		for rows.Next() {
			var address gin.H
			var id int
			var addressLine1, city, country, postalCode, addressType string
			var isDefault bool
			if err := rows.Scan(&id, &addressLine1, &city, &country, &postalCode, &addressType, &isDefault); err != nil {
				log.Println("Error scanning address row:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}
			address = gin.H{
				"id":            id,
				"address_line1": addressLine1,
				"city":          city,
				"country":       country,
				"postal_code":   postalCode,
				"type":          addressType,
				"is_default":    isDefault,
			}
			addresses = append(addresses, address)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Billing addresses retrieved",
			"addresses": addresses,
		})
	}
}

func GetSingleBillingAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		addressIDStr := c.Param("address_id")
		addressID, err := strconv.Atoi(addressIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		var addressLine1, city, country, postalCode, addressType string
		var isDefault bool
		err = db.QueryRow(`
			SELECT address_line1, city, country, postal_code, type, is_default
			FROM billing_addresses
			WHERE id = $1 AND user_id = $2
		`, addressID, userID).Scan(&addressLine1, &city, &country, &postalCode, &addressType, &isDefault)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Billing address not found"})
			return
		} else if err != nil {
			log.Println("Error fetching billing address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"address": gin.H{
				"id":            addressID,
				"address_line1": addressLine1,
				"city":          city,
				"country":       country,
				"postal_code":   postalCode,
				"type":          addressType,
				"is_default":    isDefault,
			},
		})
	}
}

func UpdateBillingAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		addressIDStr := c.Param("address_id")
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

		result, err := db.Exec(`
			UPDATE billing_addresses
			SET address_line1 = $1, city = $2, country = $3, postal_code = $4, type = $5
			WHERE id = $6 AND user_id = $7
		`, input.AddressLine1, input.City, input.Country, input.PostalCode, input.Type, addressID, userID)
		if err != nil {
			log.Println("Error updating billing address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Billing address not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Billing address updated"})
	}
}

func DeleteBillingAddress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		addressIDStr := c.Param("address_id")
		addressID, err := strconv.Atoi(addressIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		result, err := db.Exec(`
			DELETE FROM billing_addresses
			WHERE id = $1 AND user_id = $2
		`, addressID, userID)
		if err != nil {
			log.Println("Error deleting billing address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Billing address not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Billing address deleted"})
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
