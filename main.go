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

	FipContextService, err := NewFipContextService(cfg.FIPApiURL, cfg.WikiApiURL, cfg.DbpediaURL)
	if err != nil {
		log.Fatalf("couldnt create FipContextService: %w", err)
	}

	server := gin.Default()
	SetupRoutes(server, cfg, FipContextService)

	// Error middleware
	server.Use(ErrorMiddleware)
	server.Run()
}