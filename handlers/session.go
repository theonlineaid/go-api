package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"my-api/models"

	"github.com/gin-gonic/gin"
)

func Session(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var user struct {
			ID          int
			Username    string
			Email       string
			Role        string
			PhoneNumber *string
			Image       *string
			IsVerified  bool
			IsBlocked   bool
		}
		err := db.QueryRow(`
			SELECT id, username, email, role, phone_number, image, is_verified, is_blocked
			FROM users
			WHERE id = $1
		`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Role,
			&user.PhoneNumber, &user.Image, &user.IsVerified, &user.IsBlocked)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			log.Println("Error querying user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		// Fetch addresses
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

		// Fetch shipping addresses
		shippingRows, err := db.Query(`
			SELECT id, user_id, address_line1, city, country, postal_code, type, is_default, created_at
			FROM shipping_addresses
			WHERE user_id = $1
		`, userID)
		if err != nil {
			log.Println("Error querying shipping addresses:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer shippingRows.Close()

		var shippingAddresses []models.ShippingAddress
		for shippingRows.Next() {
			var addr models.ShippingAddress
			err := shippingRows.Scan(&addr.ID, &addr.UserID, &addr.AddressLine1, &addr.City, &addr.Country,
				&addr.PostalCode, &addr.Type, &addr.IsDefault, &addr.CreatedAt)
			if err != nil {
				log.Println("Error scanning shipping address:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}
			shippingAddresses = append(shippingAddresses, addr)
		}

		// Fetch billing addresses
		billingRows, err := db.Query(`
			SELECT id, user_id, address_line1, city, country, postal_code, type, is_default, created_at
			FROM billing_addresses
			WHERE user_id = $1
		`, userID)
		if err != nil {
			log.Println("Error querying billing addresses:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer billingRows.Close()

		var billingAddresses []models.BillingAddress
		for billingRows.Next() {
			var addr models.BillingAddress
			err := billingRows.Scan(&addr.ID, &addr.UserID, &addr.AddressLine1, &addr.City, &addr.Country,
				&addr.PostalCode, &addr.Type, &addr.IsDefault, &addr.CreatedAt)
			if err != nil {
				log.Println("Error scanning billing address:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}
			billingAddresses = append(billingAddresses, addr)
		}

		// Fetch login sessions
		sessionRows, err := db.Query(`
			SELECT id, user_id, browser, os, device, ip_address, login_at
			FROM login_sessions
			WHERE user_id = $1
			ORDER BY login_at DESC
		`, userID)
		if err != nil {
			log.Println("Error querying login sessions:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer sessionRows.Close()

		var loginSessions []models.LoginSession
		for sessionRows.Next() {
			var session models.LoginSession
			err := sessionRows.Scan(&session.ID, &session.UserID, &session.Browser, &session.OS,
				&session.Device, &session.IPAddress, &session.LoginAt)
			if err != nil {
				log.Println("Error scanning login session:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}
			loginSessions = append(loginSessions, session)
		}

		response := gin.H{
			"id":                 user.ID,
			"username":           user.Username,
			"email":              user.Email,
			"role":               user.Role,
			"phone_number":       user.PhoneNumber,
			"image":              user.Image,
			"is_verified":        user.IsVerified,
			"is_blocked":         user.IsBlocked,
			"addresses":          addresses,
			"shipping_addresses": shippingAddresses,
			"billing_addresses":  billingAddresses,
			"login_sessions":     loginSessions,
		}

		c.JSON(http.StatusOK, response)
	}
}
