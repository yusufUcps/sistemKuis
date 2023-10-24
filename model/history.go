package model

import (
	"gorm.io/gorm"
)

type HistoryAnswer struct {
	gorm.Model
	User_id uint `json:"user_id"`
	Quiz_id uint   `json:"quiz_id"`
	Score int `json:"score"`
}