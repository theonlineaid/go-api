// handlers/auth.go
package handlers

import (
	"database/sql"
	"log"
	"my-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"` // e.g., user, admin
}

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Validate role
		if input.Role != "user" && input.Role != "admin" && input.Role != "manager" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be user, admin, or manager"})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		// Insert user
		query := `INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id`
		var userID int
		err = db.QueryRow(query, input.Username, string(hashedPassword), input.Role).Scan(&userID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
				return
			}
			log.Println("Error inserting user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user_id": userID})
	}
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user struct {
			Password string
			Role     string
		}
		err := db.QueryRow("SELECT password, role FROM users WHERE username = $1", input.Username).
			Scan(&user.Password, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			log.Println("Error querying user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT with role
		token, err := utils.GenerateJWT(input.Username, user.Role)
		if err != nil {
			log.Println("Error generating JWT:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
			"role":    user.Role,
		})
	}
}
