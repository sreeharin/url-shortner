package main

import (
	"log"

	"gorm.io/gorm"
)

func insertData(db *gorm.DB, url URL) error {
	record := UrlDB{URL: url}
	if err := db.Create(&record).Error; err != nil {
		log.Printf("Error inserting data: %v", err)
		return err
	}
	log.Printf("Inserted record: ID=%d, Original=%s, Shortened=%s", record.ID, record.Original, record.Shortened)
	return nil
}
