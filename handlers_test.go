package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestEnvironment() (*gin.Engine, *gorm.DB, Handler) {
	router := gin.Default()
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&UrlDB{})
	handler := Handler{db: db}

	return router, db, handler
}

// TestHandleFormInput tests the handleFormInput function.
// It checks if the function correctly handles a valid input and returns the expected response.
// It also verifies that the data is correctly inserted into the database.
func TestHandleFormInput(t *testing.T) {
	router, db, handler := setupTestEnvironment()

	convertedURL := convertURL("example.com")
	exampleInput := formInput{Url: convertedURL.Original}
	inputJson, _ := json.Marshal(exampleInput)

	router.POST("/", handler.handleFormInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(inputJson))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code: %d, got: %d", http.StatusCreated, w.Code)
	}

	var url URL
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
		var urlDB UrlDB
		res := db.First(&urlDB)

		if res.Error != nil {
			if res.Error != gorm.ErrRecordNotFound {
				t.Errorf("No record found in DB: %v", res.Error)
			} else {
				t.Error(res.Error)
			}
		} else {

			if urlDB.Original != url.Original {
				t.Errorf("Expected original URL in DB: %s, got: %s", url.Original, urlDB.Original)
			}

			if urlDB.Shortened != convertedURL.Shortened {
				t.Errorf("Expected shortened URL in DB: %s, got: %s", convertedURL.Shortened, urlDB.Shortened)
			}

			if !strings.HasPrefix(urlDB.Original, "http://") {
				t.Errorf("Expected original URL to start with http://, got: %s", urlDB.Original)
			}

		}
	})

}

// TestHandleQuery tests the handleQuery function.
// It checks if the function correctly handles a valid query and returns the expected response.
// It also verifies the user is redirected to the original URL.
func TestHandleQuery(t *testing.T) {
	router, db, handler := setupTestEnvironment()

	router.GET("/", handler.handleQuery)
	convertedURL := convertURL("example.com")
	insertData(db, convertedURL)

	t.Run("TestHandleQueryValid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/?url=%s", convertedURL.Shortened), nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusMovedPermanently {
			t.Errorf("Expected status code: %d, got: %d", http.StatusMovedPermanently, w.Code)
		}

		if w.Header().Get("Location") != convertedURL.Original {
			t.Errorf("Expected Location header: %s, got: %s", convertedURL.Original, w.Header().Get("Location"))
		}
	})

	t.Run("TestHandleQueryNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?url=notfound.com", nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code: %d, got: %d", http.StatusNotFound, w.Code)
		}
	})
}
