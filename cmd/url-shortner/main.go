package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sreeharin/url-shortner/internal/handlers"
	"github.com/sreeharin/url-shortner/internal/models"
)

func main() {
	dsn := "host=db user=user password=password dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.UrlDB{})

	handler := handlers.Handler{DB: db}

	router := gin.Default()
	router.POST("/", handler.HandleFormInput)
	router.GET("/", handler.HandleQuery)
	router.Run()
}
