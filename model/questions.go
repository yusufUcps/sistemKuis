package model

import (
	"gorm.io/gorm"
)

type Questions struct {
	gorm.Model
	Quiz_id uint `json:"id_quiz" gorm:"type:varchar(255)"`
	Question    string `json:"question" gorm:"type:text"`
	Options []Options `json:"options" gorm:"foreignKey:Quiz_id"`
}