package model

import (
	"gorm.io/gorm"
)

type HistoryAnswer struct {
	gorm.Model
	User_id uint `json:"user_id"`
	Name string `json:"name" gorm:"type:varchar(255)"`
	Quiz_id uint   `json:"quiz_id"`
	Title string `json:"title" gorm:"type:varchar(225)"`
	Score int `json:"score"`
}