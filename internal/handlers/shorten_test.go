package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gorm.io/gorm"

	"github.com/sreeharin/url-shortner/internal/models"
	"github.com/sreeharin/url-shortner/internal/utils"
)

func TestShortenURL(t *testing.T) {
	router, DB := SetupTestEnvironment(t)
	handler := Handler{DB: DB}

	convertedURL := utils.ConvertURL("example.com")
	exampleInput := formInput{Url: convertedURL.Original}
	inputJson, _ := json.Marshal(exampleInput)

	router.POST("/", handler.ShortenURL)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(inputJson))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code: %d, got: %d", http.StatusCreated, w.Code)
	}

	var url models.URL
	json.Unmarshal(w.Body.Bytes(), &url)

	t.Run("TestHandleFormInputValid", func(t *testing.T) {
		if url.Original != exampleInput.Url {
			t.Errorf("Expected original URL: %s, got: %s", exampleInput.Url, url.Original)
		}

		if url.Shortened != convertedURL.Shortened {
			t.Errorf("Expected shortened URL: %s, got: %s", convertedURL.Shortened, url.Shortened)
		}

	})

	t.Run("TestDBInsertion", func(t *testing.T) {
		var url models.URL
		res := DB.First(&url)

		if res.Error != nil {
			if res.Error != gorm.ErrRecordNotFound {
				t.Errorf("No record found in DB: %v", res.Error)
			} else {
				t.Error(res.Error)
			}
		} else {

			if url.Original != url.Original {
				t.Errorf("Expected original URL in DB: %s, got: %s", url.Original, url.Original)
			}

			if url.Shortened != convertedURL.Shortened {
				t.Errorf("Expected shortened URL in DB: %s, got: %s", convertedURL.Shortened, url.Shortened)
			}

			if !strings.HasPrefix(url.Original, "http://") {
				t.Errorf("Expected original URL to start with http://, got: %s", url.Original)
			}

		}
	})

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
