package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"my-api/utils" // Update this path to match your project structure

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register handler
func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Bind incoming JSON to the user struct
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Hash the password before storing
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		// Prepare SQL query to insert the new user
		query := `INSERT INTO users (username, password) VALUES ($1, $2)`
		_, err = db.Exec(query, user.Username, string(hashedPassword))
		if err != nil {
			log.Println("Error inserting user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = $1", input.Username).Scan(&storedPassword)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, _ := utils.GenerateJWT(input.Username)
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
