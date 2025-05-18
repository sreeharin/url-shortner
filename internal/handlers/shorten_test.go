package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sreeharin/url-shortner/internal/models"
	"github.com/sreeharin/url-shortner/internal/utils"
)

func validateShortenURL(t *testing.T, want, got string) {
	t.Helper()

	if want != got {
		t.Errorf("Failed to shorten URL. Want :%s Got :%s", want, got)
	}

}

func TestShortenURL(t *testing.T) {
	router, DB := SetupTestEnvironment(t)
	handler := Handler{DB: DB}

	router.POST("/", handler.ShortenURL)

	var (
		w     *httptest.ResponseRecorder
		req   *http.Request
		input []byte
	)

	tt := []struct {
		original   string
		want       string
		statusCode int
	}{
		{
			original:   "example.com",
			want:       utils.ConvertURL("example.com").Shortened,
			statusCode: http.StatusCreated,
		},

		// Adding http:// prefix should also produce the same shortened code of example.com
		{
			original:   "http://example.com",
			want:       utils.ConvertURL("example.com").Shortened,
			statusCode: http.StatusOK,
		},

		{
			original:   "http://google.com",
			want:       utils.ConvertURL("http://google.com").Shortened,
			statusCode: http.StatusCreated,
		},
		{
			original:   "http://amazon.com/login",
			want:       utils.ConvertURL("http://amazon.com").Shortened,
			statusCode: http.StatusCreated,
		},
	}

	for _, testCase := range tt {
		input, _ = json.Marshal(formInput{Url: testCase.original})

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/", bytes.NewBuffer(input))

		router.ServeHTTP(w, req)

		if w.Code != testCase.statusCode {
			t.Errorf("Expected status code: %d, got: %d", testCase.statusCode, w.Code)
		}
		var url models.URL
		json.Unmarshal(w.Body.Bytes(), &url)

		validateShortenURL(t, testCase.want, url.Shortened)
	}

}

func TestRedirectURL(t *testing.T) {
	router, DB := SetupTestEnvironment(t)
	handler := Handler{DB: DB}

	router.GET("/:url", handler.RedirectURL)
	convertedURL := utils.ConvertURL("example.com")

	DB.Create(&convertedURL)

	t.Run("TestParamValid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", convertedURL.Shortened), nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusMovedPermanently {
			t.Errorf("Expected status code: %d, got: %d", http.StatusMovedPermanently, w.Code)
		}

		if w.Header().Get("Location") != convertedURL.Original {
			t.Errorf("Expected Location header: %s, got: %s", convertedURL.Original, w.Header().Get("Location"))
		}
	})

	t.Run("TestQueryNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/notfound.com", nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code: %d, got: %d", http.StatusNotFound, w.Code)
		}
	})
}
