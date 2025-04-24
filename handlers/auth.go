// handlers/auth.go
package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"my-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/mssola/useragent"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username    string  `json:"username" binding:"required,min=3"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=8"`
	Role        string  `json:"role" binding:"required,oneof=user admin manager"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Image       *string `json:"image,omitempty"`
}

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		query := `INSERT INTO users (username, email, password, role, phone_number, image, is_verified, is_blocked, created_at, updated_at)
		          VALUES ($1, $2, $3, $4, $5, $6, FALSE, FALSE, NOW(), NOW()) RETURNING id`
		var userID int
		err = db.QueryRow(query, input.Username, input.Email, string(hashedPassword), input.Role, input.PhoneNumber, input.Image).Scan(&userID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				if strings.Contains(pqErr.Detail, "username") {
					c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
				} else {
					c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
				}
				return
			}
			log.Println("Error inserting user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		// accessToken, err := utils.GenerateAccessToken(uint(userID), input.Username, input.Email, string(hashedPassword), input.Role)
		// if err != nil {
		// 	log.Println("Error generating access token:", err)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		// 	return
		// }
		// refreshToken, err := utils.GenerateRefreshToken(uint(userID), input.Username, input.Email, string(hashedPassword), input.Role)
		// if err != nil {
		// 	log.Println("Error generating refresh token:", err)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		// 	return
		// }

		c.JSON(http.StatusOK, gin.H{
			"message": "User created",
			"user_id": userID,
			// "access_token":  accessToken,
			// "refresh_token": refreshToken,
			"role": input.Role,
		})
	}
}

type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user struct {
			ID        int
			Username  string
			Email     string
			Password  string
			Role      string
			IsBlocked bool
		}
		err := db.QueryRow(`
			SELECT id, username, email, password, role, is_blocked
			FROM users
			WHERE username = $1 OR email = $1
		`, input.Identifier).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.IsBlocked)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			log.Println("Error querying user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		if user.IsBlocked {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account is blocked"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		ua := useragent.New(c.GetHeader("User-Agent"))
		browser, browserVersion := ua.Browser()
		browserInfo := browser + " " + browserVersion
		osInfo := ua.OS()
		deviceInfo := ua.Model()
		if deviceInfo == "" {
			deviceInfo = "Unknown"
		}
		ipAddress := c.ClientIP()

		_, err = db.Exec(`
			INSERT INTO login_sessions (user_id, browser, os, device, ip_address, login_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
		`, user.ID, browserInfo, osInfo, deviceInfo, ipAddress)
		if err != nil {
			log.Println("Error storing login session:", err)
		}

		accessToken, err := utils.GenerateAccessToken(uint(user.ID), user.Username, user.Email, user.Password, user.Role)
		if err != nil {
			log.Println("Error generating access token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(uint(user.ID), user.Username, user.Email, user.Password, user.Role)
		if err != nil {
			log.Println("Error generating refresh token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "Login successful",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"role":          user.Role,
		})
	}
}

func RefreshToken(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		claims, err := utils.ParseToken(input.RefreshToken)
		if err != nil || claims.Type != "refresh" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		var user struct {
			ID       int
			Username string
			Email    string
			Password string
			Role     string
		}
		err = db.QueryRow(`
			SELECT id, username, email, password, role
			FROM users
			WHERE id = $1
		`, claims.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			log.Println("Error querying user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		accessToken, err := utils.GenerateAccessToken(uint(user.ID), user.Username, user.Email, user.Password, user.Role)
		if err != nil {
			log.Println("Error generating access token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Token refreshed",
			"access_token": accessToken,
			"role":         user.Role,
		})
	}
}

func Logout(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		_, err := db.Exec(`
			DELETE FROM login_sessions
			WHERE user_id = $1
		`, userID)
		if err != nil {
			log.Println("Error deleting login sessions:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Logged out successfully, all sessions cleared",
		})
	}
}
