package model

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(255)"`
	Email    string `json:"email" gorm:"type:varchar(20);unique"`
	Password string `json:"password" `
}

