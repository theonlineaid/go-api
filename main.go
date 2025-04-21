package main

import (
	"database/sql"
	"log"
	"my-api/handlers" // Update the import path according to your project structure

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	connStr := "user=admin password=secret dbname=mydb host=db port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Test the DB connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	// Setup Gin router
	router := gin.Default()

	// Register and Login routes
	router.POST("/register", handlers.Register(db))
	router.POST("/login", handlers.Login(db))
	router.POST("/product", handlers.AddProduct(db))
	router.GET("/product/:id", handlers.GetProductByID(db))

	// Start the server
	router.Run("0.0.0.0:8080")
}
