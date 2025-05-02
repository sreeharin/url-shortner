package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHandleFormInput(t *testing.T) {
	router := setupRouter()

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&UrlDB{})

	handler := handler{db: db}
	router.POST("/", handler.handleFormInput)

	w := httptest.NewRecorder()

	convertedURL := convertURL("http://example.com")
	exampleInput := formInput{Url: convertedURL.Original}
	inputJson, _ := json.Marshal(exampleInput)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(inputJson))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code: %d, got: %d", http.StatusCreated, w.Code)
	}

	var url URL
	json.Unmarshal(w.Body.Bytes(), &url)

	if url.Original != exampleInput.Url {
		t.Errorf("Expected original URL: %s, got: %s", exampleInput.Url, url.Original)
	}

	if url.Shortened != convertedURL.Shortened {
		t.Errorf("Expected shortened URL: %s, got: %s", convertedURL.Shortened, url.Shortened)
	}

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

	}

}
