package models

import "gorm.io/gorm"

type URL struct {
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
}

type UrlDB struct {
	ID uint `gorm:"primaryKey"`
	URL
}

type User struct {
	*gorm.Model `json:"-"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
}
