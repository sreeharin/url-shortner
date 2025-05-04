package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/sreeharin/url-shortner/internal/models"
)

func SetupTestEnvironment(t *testing.T) (router *gin.Engine, db *gorm.DB) {
	t.Helper()

	router = setupRouter()
	db = setupDB()
	db.AutoMigrate(&models.UrlDB{})

	return
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func setupDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	return db
}
