package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type formInput struct {
	Url string `json:"url" binding:"required"`
}

func handleFormInput(c *gin.Context) {
	var input formInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is missing"})
		return
	}

	converted := convertURL(input.Url)

	c.JSON(http.StatusOK, converted)
}

func handleQuery(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is missing"})
		return
	}

	c.JSON(http.StatusOK, url)
}
