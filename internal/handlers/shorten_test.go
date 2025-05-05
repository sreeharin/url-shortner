package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
		var urlDB models.URL
		DB.First(&urlDB)

		if urlDB.Original != url.Original {
			t.Errorf("Expected original URL in DB: %s, got: %s", url.Original, urlDB.Original)
		}

		if urlDB.Shortened != convertedURL.Shortened {
			t.Errorf("Expected shortened URL in DB: %s, got: %s", convertedURL.Shortened, urlDB.Shortened)
	}

		if !strings.HasPrefix(urlDB.Original, "http://") {
			t.Errorf("Expected original URL to start with http://, got: %s", urlDB.Original)
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
