package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type formInput struct {
	Url string `json:"url" binding:"required"`
}

type handler struct {
	db *gorm.DB
}

// handleFormInput handles the form input from the client.
// It expects a JSON body with a "url" field.
// It converts the URL to a shortened version and inserts it into the database.
func (h handler) handleFormInput(c *gin.Context) {
	var input formInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is missing"})
		return
	}

	converted := convertURL(input.Url)

	if err := insertData(h.db, converted); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert data"})
		return
	}

	c.JSON(http.StatusOK, converted)
}

func (h handler) handleQuery(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is missing"})
		return
	}

	c.JSON(http.StatusOK, url)
}
