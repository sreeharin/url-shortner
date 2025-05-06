package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuth(t *testing.T) {
	router := gin.New()
	router.Use(Auth())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	t.Run("Status 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		router.ServeHTTP(w, req)

		log.Println(w.Body.String())

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code: %d got: %d", http.StatusUnauthorized, w.Code)
		}
	})

}
