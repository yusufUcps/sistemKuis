package model

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	User_id uint `json:"user_id"`
	Title     string `json:"title" gorm:"type:varchar(225)"`
	Description    string `json:"description" gorm:"type:text"`
	Start_date time.Time `json:"start_date" gorm:"type:date"`
	End_date time.Time `json:"end_date" gorm:"type:date"`
	Questions []Questions `json:"questions" gorm:"foreignKey:Quiz_id"`
}