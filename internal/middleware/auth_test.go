package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sreeharin/url-shortner/internal/models"
)

func TestAuth(t *testing.T) {
	router := gin.New()

	router.GET("/", Auth(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	t.Run("Status 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code: %d got: %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Authorized", func(t *testing.T) {
		claims := models.CustomClaim{
			User: 1,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("secret"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code: %d got: %d", http.StatusOK, w.Code)
		}
	})

}
