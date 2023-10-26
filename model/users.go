package model

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
    Name     string `json:"name" gorm:"type:varchar(255);not null"`
    Email    string `json:"email" gorm:"type:varchar(20);unique;not null"`
    Password string `json:"password" gorm:"not null"`

}

