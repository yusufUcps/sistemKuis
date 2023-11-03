package model

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
    Name     string `json:"name" gorm:"type:varchar(255);not null"`
    Email    string `json:"email" gorm:"type:varchar(50);unique;not null"`
    Password string `json:"password" gorm:"not null"`
}

type UserRegisterRes struct {
	Id uint  `json:"id" `
    Created_at time.Time`json:"created_at" `
    Updatad_at time.Time`json:"updated_at" `
    Name     string `json:"name" `
    Email    string `json:"email" `
}

type MyProfileRes struct {
	Id uint  `json:"id" `
    Created_at time.Time`json:"created_at" `
    Updatad_at time.Time`json:"updated_at" `
    Name     string `json:"name" `
    Email    string `json:"email" `
}

func ConvertRegisterRes(user *Users) *UserRegisterRes {
    registerRes := UserRegisterRes{
        Id:         user.ID,
        Created_at: user.CreatedAt,
        Updatad_at: user.UpdatedAt,
        Name:       user.Name,
        Email:      user.Email,
    }
    return &registerRes
}

func ConvertMyProfileRes(user *Users) *MyProfileRes {
    myProfileRes := MyProfileRes{
        Id:         user.ID,
        Created_at: user.CreatedAt,
        Updatad_at: user.UpdatedAt,
        Name:       user.Name,
        Email:      user.Email,
    }
    return &myProfileRes
}

