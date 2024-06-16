package main

import (
	"log"
	"sumup/asset/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Define the routes
	r.POST("/create", createAccountHandler)
	r.POST("/deposit", depositHandler)
	r.POST("/transfer", transferHandler)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func createAccountHandler(c *gin.Context) {
	// Your handler logic
}

func depositHandler(c *gin.Context) {
	// Your handler logic
}

func transferHandler(c *gin.Context) {
	// Your handler logic
}
