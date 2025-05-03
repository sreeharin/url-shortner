package db

import (
	"gorm.io/gorm"

	"github.com/sreeharin/url-shortner/internal/models"
)

func InsertData(db *gorm.DB, url models.URL) error {
	record := models.UrlDB{URL: url}
	if err := db.Create(&record).Error; err != nil {
		return err
	}
	return nil
}
