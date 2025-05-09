package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *gin.Engine, cfg *Config, fipContextService FipContextService) {

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	server.GET("/current", func(ctx *gin.Context) {
		from := ctx.Query("from")

		song, err := fipContextService.Current(from)

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		} else {
			ctx.JSON(http.StatusOK, song)
		}
	})
}
