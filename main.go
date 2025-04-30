package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/", handleFormInput)
	router.GET("/", handleQuery)
	router.Run()
}
