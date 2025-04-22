// main.go
package main

import (
	"database/sql"
	"log"
	"my-api/handlers"
	"os"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"}, // your frontend domains
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.POST("/register", handlers.Register(db))
	r.POST("/login", handlers.Login(db))
	r.GET("/products", handlers.GetAllProducts(db))
	r.GET("/products/:id", handlers.GetProductByID(db))
	r.GET("/brands", handlers.GetAllBrands(db))
	r.GET("/brands/:id", handlers.GetBrandByID(db))
	r.GET("/categories", handlers.GetAllCategories(db))
	r.GET("/categories/:id", handlers.GetCategoryByID(db))
	r.GET("/subcategories", handlers.GetAllSubCategories(db))
	r.GET("/subcategories/:id", handlers.GetSubCategoryByID(db))

	// Protected routes (admin-only for product, brand, etc.)
	adminGroup := r.Group("/admin")
	adminGroup.Use(handlers.JWTAuthMiddleware("admin"))
	{
		adminGroup.POST("/products", handlers.AddProduct(db))
		adminGroup.GET("/products", handlers.GetAllProducts(db))
		adminGroup.GET("/products/:id", handlers.GetProductByID(db))
		adminGroup.POST("/brands", handlers.CreateBrand(db))
		adminGroup.PUT("/brands/:id", handlers.UpdateBrand(db))
		adminGroup.DELETE("/brands/:id", handlers.DeleteBrand(db))
		adminGroup.POST("/categories", handlers.CreateCategory(db))
		adminGroup.PUT("/categories/:id", handlers.UpdateCategory(db))
		adminGroup.DELETE("/categories/:id", handlers.DeleteCategory(db))
		adminGroup.GET("/categories/:id", handlers.GetCategoryByID(db))
		adminGroup.POST("/subcategories", handlers.CreateSubCategory(db))
		adminGroup.PUT("/subcategories/:id", handlers.UpdateSubCategory(db))
		adminGroup.DELETE("/subcategories/:id", handlers.DeleteSubCategory(db))
		adminGroup.GET("/subcategories", handlers.GetAllSubCategories(db))
		adminGroup.GET("/subcategories/:id", handlers.GetSubCategoryByID(db))
		adminGroup.POST("/attribute-values", handlers.CreateAttributeValue(db))
	}

	// User routes (example: users can view products)
	userGroup := r.Group("/user")
	userGroup.Use(handlers.JWTAuthMiddleware("user"))
	{
		userGroup.GET("/products/:id", handlers.GetProductByID(db))
		userGroup.GET("/products", handlers.GetAllProducts(db))
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
