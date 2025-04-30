package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=db user=user password=password dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&UrlDB{})

	handler := handler{db: db}

	router := gin.Default()
	router.POST("/", handler.handleFormInput)
	router.GET("/", handler.handleQuery)
	router.Run()
}
