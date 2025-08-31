package main

import (
	"log"

	"stokq-backend/config"
	"stokq-backend/initializers"
	"stokq-backend/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	// Load environment variables
	initializers.LoadEnvVariables()

	// Connect to database
	config.ConnectDatabase()
}

func main() {
	// Set Gin mode (release for production)
	gin.SetMode(gin.DebugMode)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Get port from environment
	port := initializers.GetEnv("PORT", "8080")

	log.Printf("🚀 StokQ API Server starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/health", port)
	log.Printf("📝 API Documentation: http://localhost:%s/api/v1", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
