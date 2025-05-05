package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sreeharin/url-shortner/internal/handlers"
	"github.com/sreeharin/url-shortner/internal/middleware"
	"github.com/sreeharin/url-shortner/internal/models"
)

func main() {
	dsn := "host=db user=user password=password dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.URL{})
	db.AutoMigrate(&models.User{})

	handler := handlers.Handler{DB: db}
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))

	router.POST("/", handler.ShortenURL)
	router.GET("/:url", handler.RedirectURL)

	router.Run()
}
