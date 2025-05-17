package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sreeharin/url-shortner/internal/metrics"
	"github.com/sreeharin/url-shortner/internal/models"
	"github.com/sreeharin/url-shortner/internal/utils"
)

type formInput struct {
	Url string `json:"url" binding:"required"`
}

// ShortenURL handles the URL shortening process.
// It expects a JSON payload with the original URL.
// It converts the original URL to a shortened version and stores it in the database.
func (h *Handler) ShortenURL(c *gin.Context) {
	var input formInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed form data"})
		return
	}

	if !strings.HasPrefix(input.Url, "http") {
		input.Url = "http://" + input.Url
	}

	var urlDB models.URL
	if err := h.DB.Where("original = ?", input.Url).First(&urlDB).Error; err == nil {
		c.JSON(http.StatusOK, urlDB)
		return
	}

	converted := utils.ConvertURL(input.Url)

	if err := h.DB.Create(&converted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		return
	}

	metrics.URLShortenRequests.Inc()
	c.JSON(http.StatusCreated, converted)
}

// RedirectURL handles the redirection from the shortened URL to the original URL.
// It expects the shortened URL as a URL parameter.
func (h *Handler) RedirectURL(c *gin.Context) {
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Url is missing"})
		return
	}

	var urlDB models.URL
	if err := h.DB.Where("shortened = ?", url).First(&urlDB).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Url not found"})
		return
	}

	metrics.URLRedirectRequests.WithLabelValues(url).Inc()
	c.Redirect(http.StatusMovedPermanently, urlDB.Original)
}
