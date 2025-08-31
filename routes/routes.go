package routes

import (
	"stokq-backend/controllers"
	"stokq-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// CORS middleware - Add this for cross-origin requests
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "StokQ API is running",
		})
	})

	// API version 1 group
	api := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Protected routes (authentication required)
	protected := api.Group("/")
	protected.Use(middleware.RequireAuth)
	{
		// Product routes
		products := protected.Group("/products")
		{
			products.POST("/", controllers.CreateProduct)
			products.GET("/", controllers.GetProducts)
			products.GET("/:id", controllers.GetProductByID)
			products.PUT("/:id", controllers.UpdateProduct)
			products.DELETE("/:id", controllers.DeleteProduct)
		}

		// Stock routes
		stock := protected.Group("/stock")
		{
			stock.POST("/in", controllers.StockIn)
			stock.POST("/out", controllers.StockOut)
		}
	}
}
