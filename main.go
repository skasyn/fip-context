package main

import "github.com/gin-gonic/gin"

func main() {
	server := gin.Default()
	SetupRoutes(server)

	// Error middleware
	server.Use(ErrorMiddleware)
	server.Run()
}