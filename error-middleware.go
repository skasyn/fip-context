package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		log.Printf(err.Error())
	}

	// status -1 doesn't overwrite existing status code
	c.JSON(-1, "")
}