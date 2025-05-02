package main

import (
	"gorm.io/gorm"
)

func insertData(db *gorm.DB, url URL) error {
	record := UrlDB{URL: url}
	if err := db.Create(&record).Error; err != nil {
		return err
	}
	return nil
}
