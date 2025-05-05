package models

import "gorm.io/gorm"

type URL struct {
	*gorm.Model `json:"-"`
	Original    string `json:"original"`
	Shortened   string `json:"shortened"`
}

type User struct {
	*gorm.Model `json:"-"`
	Username    string `json:"username" binding:"required" gorm:"unique"`
	Password    string `json:"password" binding:"required"`
}
