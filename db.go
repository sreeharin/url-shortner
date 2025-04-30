package main

import (
	"gorm.io/gorm"
)

func insertData(db *gorm.DB, url URL) error {
	if err := db.Create(&UrlDB{URL: url}).Error; err != nil {
		return err
	}
	return nil
}
