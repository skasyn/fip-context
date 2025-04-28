package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/skasyn/fip-context/fipcontextrepo"
)

func main() {
	cfg, err := LoadConfig()

	if err != nil {
		log.Fatalf("error while loading config: %s", err.Error())
	}

	dbConnect, err := fipcontextrepo.Connect(cfg.psqlConnStr)
	if err != nil {
		log.Fatalf("couldnt connect to db: %w", err)
	}
	
	fipContextRepo := fipcontextrepo.NewFipSongRepository(dbConnect)

	FipContextService, err := NewFipContextService(cfg, fipContextRepo)
	if err != nil {
		log.Fatalf("couldnt create FipContextService: %w", err)
	}

	server := gin.Default()
	SetupRoutes(server, cfg, FipContextService)

	// Error middleware
	server.Use(ErrorMiddleware)
	server.Run()
}