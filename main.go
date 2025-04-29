package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jxskiss/base62"
)

func convertURL(original string) URL {
	var converted URL
	converted.Original = original

	LIMIT := 6
	encoded := base62.EncodeToString([]byte(original))

	if len(encoded) > LIMIT {
		converted.Shortened = encoded[len(encoded)-LIMIT:]
	} else {
		converted.Shortened = encoded
	}

	return converted
}

func handler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "URL is missing"})
		return
	}
	c.JSON(http.StatusOK, convertURL(url))
}

func main() {
	router := gin.Default()
	router.GET("/", handler)
	router.Run()
}
