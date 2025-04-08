package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *gin.Engine) {
	fipService := DefaultFipService{}

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	server.GET("/current", func(ctx *gin.Context) {
		from := ctx.Query("from")

		song, err := fipService.GetCurrentSong(from)

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, song)
	})
}
