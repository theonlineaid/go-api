// main.go
package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"my-api/handlers"
	"os"
)

func main() {
	// Database connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Public routes
	r.POST("/register", handlers.Register(db))
	r.POST("/login", handlers.Login(db))

	// Protected routes (admin-only for product, brand, etc.)
	adminGroup := r.Group("/admin")
	adminGroup.Use(handlers.JWTAuthMiddleware("admin"))
	{
		adminGroup.POST("/products", handlers.AddProduct(db))
		adminGroup.GET("/products/:id", handlers.GetProductByID(db))
		adminGroup.POST("/brands", handlers.CreateBrand(db))
		adminGroup.POST("/categories", handlers.CreateCategory(db))
		adminGroup.POST("/subcategories", handlers.CreateSubcategory(db))
		adminGroup.POST("/attributes", handlers.CreateAttribute(db))
		adminGroup.POST("/attribute-values", handlers.CreateAttributeValue(db))
	}

	// User routes (example: users can view products)
	userGroup := r.Group("/user")
	userGroup.Use(handlers.JWTAuthMiddleware("user"))
	{
		userGroup.GET("/products/:id", handlers.GetProductByID(db))
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
