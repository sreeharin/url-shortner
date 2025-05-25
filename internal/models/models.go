package models

import (
	"github.com/sreeharin/url-shortner/internal/utils"
	"gorm.io/gorm"
)

type URL struct {
	*gorm.Model `json:"-"`
	Original    string `json:"original" gorm:"unique"`
	Shortened   string `json:"shortened"`
}

type User struct {
	*gorm.Model `json:"-"`
	Username    string `json:"username" binding:"required" gorm:"unique"`
	Password    string `json:"password" binding:"required"`
}

func (url *URL) AfterCreate(tx *gorm.DB) (err error) {
	shortened := utils.ConvertID(url.ID)
	url.Shortened = shortened
	err = tx.Model(url).Update("Shortened", shortened).Error
	return
}
