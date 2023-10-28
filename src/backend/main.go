package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.String(http.StatusOK, "Up and running")
}

func main() {
	router := gin.Default()
	router.GET("/health", health)

	router.Run(":8000")
}
