package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg, err := LoadConfig()

	if err != nil {
		log.Fatalf("error while loading config: %v", err.Error())
	}

	server := gin.Default()
	SetupRoutes(server, cfg)

	// Error middleware
	server.Use(ErrorMiddleware)
	server.Run()
}