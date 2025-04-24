package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	_ "time"

	"my-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserInput struct {
	Username    string  `json:"username" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	Role        string  `json:"role" binding:"required,oneof=user admin"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Image       *string `json:"image,omitempty"`
	IsVerified  bool    `json:"is_verified"`
	IsBlocked   bool    `json:"is_blocked"`
}

type AdminUserUpdateInput struct {
	Username    string  `json:"username" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Role        string  `json:"role" binding:"required,oneof=user admin"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Image       *string `json:"image,omitempty"`
	IsVerified  bool    `json:"is_verified"`
	IsBlocked   bool    `json:"is_blocked"`
}

// AdminCreateUser creates a new user (admin only)
func AdminCreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AdminUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		query := `
			INSERT INTO users (username, email, role, password, phone_number, image, is_verified, is_blocked, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			RETURNING id
		`
		var userID int
		err = db.QueryRow(query, input.Username, input.Email, input.Role, hashedPassword,
			input.PhoneNumber, input.Image, input.IsVerified, input.IsBlocked).Scan(&userID)
		if err != nil {
			log.Println("Error inserting user:", err)
			if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` ||
				err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User created",
			"user_id": userID,
		})
	}
}

// AdminGetAllUsers retrieves all users (admin only)
func AdminGetAllUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT id, username, email, role, phone_number, image, is_verified, is_blocked, created_at, updated_at
			FROM users
			ORDER BY id
		`)
		if err != nil {
			log.Println("Error querying users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role,
				&user.PhoneNumber, &user.Image, &user.IsVerified, &user.IsBlocked,
				&user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				log.Println("Error scanning user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Users retrieved",
			"users":   users,
		})
	}
}

// AdminGetUser retrieves a single user by ID (admin only)
func AdminGetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user models.User
		err = db.QueryRow(`
			SELECT id, username, email, role, phone_number, image, is_verified, is_blocked, created_at, updated_at
			FROM users
			WHERE id = $1
		`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Role,
			&user.PhoneNumber, &user.Image, &user.IsVerified, &user.IsBlocked,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			log.Println("Error querying user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User retrieved",
			"user":    user,
		})
	}
}

// AdminUpdateUser updates a user by ID (admin only)
func AdminUpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var input AdminUserUpdateInput
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

		// Verify user exists
		var existsUser bool
		err = tx.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)
		`, userID).Scan(&existsUser)
		if err != nil {
			log.Println("Error verifying user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		if !existsUser {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		result, err := tx.Exec(`
			UPDATE users
			SET username = $1, email = $2, role = $3, phone_number = $4, image = $5,
				is_verified = $6, is_blocked = $7, updated_at = NOW()
			WHERE id = $8
		`, input.Username, input.Email, input.Role, input.PhoneNumber, input.Image,
			input.IsVerified, input.IsBlocked, userID)
		if err != nil {
			log.Println("Error updating user:", err)
			if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` ||
				err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User updated",
			"user_id": userID,
		})
	}
}

// AdminDeleteUser deletes a user by ID (admin only)
func AdminDeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		defer tx.Rollback()

		// Delete associated login_sessions
		_, err = tx.Exec(`
			DELETE FROM login_sessions WHERE user_id = $1
		`, userID)
		if err != nil {
			log.Println("Error deleting login sessions:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		// Delete associated shipping_addresses
		_, err = tx.Exec(`
			DELETE FROM shipping_addresses WHERE user_id = $1
		`, userID)
		if err != nil {
			log.Println("Error deleting shipping addresses:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		// Delete associated billing_addresses
		_, err = tx.Exec(`
			DELETE FROM billing_addresses WHERE user_id = $1
		`, userID)
		if err != nil {
			log.Println("Error deleting billing addresses:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		// Delete user (addresses table handled by ON DELETE CASCADE)
		result, err := tx.Exec(`
			DELETE FROM users WHERE id = $1
		`, userID)
		if err != nil {
			log.Println("Error deleting user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User deleted",
			"user_id": userID,
		})
	}
}
