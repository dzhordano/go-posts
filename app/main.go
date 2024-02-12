package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello :)")
}

func main() {
	router := gin.Default()

	router.GET("/", hello)

	router.Run()
}
