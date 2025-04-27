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
	r.POST("/refresh", handlers.RefreshToken(db))
	r.GET("/products", handlers.GetAllProducts(db))
	r.GET("/products/:id", handlers.GetProductByID(db))
	r.GET("/brands", handlers.GetAllBrands(db))
	r.GET("/brands/:id", handlers.GetBrandByID(db))
	r.GET("/categories", handlers.GetAllCategories(db))
	r.GET("/categories/:id", handlers.GetCategoryByID(db))
	r.GET("/subcategories", handlers.GetAllSubCategories(db))
	r.GET("/subcategories/:id", handlers.GetSubCategoryByID(db))

	// Protected routes (admin-only)
	adminGroup := r.Group("/admin")
	adminGroup.Use(handlers.JWTAuthMiddleware("admin"))
	{
		adminGroup.POST("/products", handlers.AddProduct(db))

		adminGroup.POST("/attributes", handlers.CreateAttribute(db))
		adminGroup.GET("/attributes", handlers.GetAllAttributes(db))
		adminGroup.GET("/attributes/:id", handlers.GetAttributeByID(db))
		adminGroup.PUT("/attributes/:id", handlers.UpdateAttribute(db))
		adminGroup.DELETE("/attributes/:id", handlers.DeleteAttribute(db))

		adminGroup.POST("/attribute-values", handlers.CreateAttributeValue(db))
		adminGroup.GET("/attribute-values/:attribute_id", handlers.GetAllAttributeValues(db))
		adminGroup.GET("/attribute-value/:id", handlers.GetAttributeValue(db))
		adminGroup.PUT("/attribute-value/:id", handlers.UpdateAttributeValue(db))
		adminGroup.DELETE("/attribute-value/:id", handlers.DeleteAttributeValue(db))

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
		adminGroup.GET("/session", handlers.Session(db))
		adminGroup.POST("/address", handlers.AdminCreateAddress(db))
		adminGroup.PUT("/address/:id", handlers.AdminUpdateAddress(db))

		adminGroup.POST("/users", handlers.AdminCreateUser(db))
		adminGroup.GET("/users", handlers.AdminGetAllUsers(db))
		adminGroup.GET("/users/:id", handlers.AdminGetUser(db))
		adminGroup.PUT("/users/:id", handlers.AdminUpdateUser(db))
		adminGroup.DELETE("/users/:id", handlers.AdminDeleteUser(db))
	}

	// User routes
	userGroup := r.Group("/user")
	userGroup.Use(handlers.JWTAuthMiddleware("user"))
	{
		userGroup.GET("/products/:id", handlers.GetProductByID(db))
		userGroup.GET("/products", handlers.GetAllProducts(db))
		userGroup.GET("/session", handlers.Session(db))
		userGroup.GET("/me", handlers.GetUserDetails(db))
		userGroup.POST("/logout", handlers.Logout(db))

		userGroup.POST("/address", handlers.AddAddress(db))
		userGroup.GET("/addresses", handlers.GetAddresses(db))
		userGroup.PUT("/address/:id", handlers.UpdateAddress(db))
		userGroup.GET("/address/:id", handlers.GetAddressByID(db))
		userGroup.DELETE("/address/:id", handlers.DeleteAddress(db))

		userGroup.POST("/shipping-address", handlers.AddShippingAddress(db))
		userGroup.GET("/shipping-addresses", handlers.GetShippingAddresses(db))
		userGroup.GET("/shipping-address/:id", handlers.GetSingleShippingAddress(db))
		userGroup.PUT("/shipping-address/:id", handlers.UpdateShippingAddress(db))
		userGroup.DELETE("/shipping-address/:id", handlers.DeleteShippingAddress(db))
		userGroup.POST("/billing-address", handlers.AddBillingAddress(db))
		userGroup.GET("/billing-addresses", handlers.GetBillingAddresses(db))
		userGroup.GET("/billing-address/:id", handlers.GetSingleBillingAddress(db))
		userGroup.PUT("/billing-address/:id", handlers.UpdateBillingAddress(db))
		userGroup.DELETE("/billing-address/:id", handlers.DeleteBillingAddress(db))

		userGroup.PATCH("/shipping-address/:id/default", handlers.SetDefaultShippingAddress(db))
		userGroup.PATCH("/billing-address/:id/default", handlers.SetDefaultBillingAddress(db))
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
