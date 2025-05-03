package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/sreeharin/url-shortner/internal/db"
	"github.com/sreeharin/url-shortner/internal/models"
	"github.com/sreeharin/url-shortner/internal/utils"
)

type formInput struct {
	Url string `json:"url" binding:"required"`
}

type Handler struct {
	DB *gorm.DB
}

// handleFormInput handles the form input from the client.
// It expects a JSON body with a "url" field.
// It converts the URL to a shortened version and inserts it into the database.
func (h *Handler) HandleFormInput(c *gin.Context) {
	var input formInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is missing"})
		return
	}

	converted := utils.ConvertURL(input.Url)

	if err := db.InsertData(h.DB, converted); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert data"})
		return
	}

	c.JSON(http.StatusCreated, converted)
}

// handleQuery handles the query for a shortened URL.
// It expects a query parameter "url" and returns the original URL if found.
func (h *Handler) HandleQuery(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is missing"})
		return
	}

	var urlDB models.UrlDB
	if err := h.DB.Where("shortened = ?", url).First(&urlDB).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, urlDB.Original)
}
