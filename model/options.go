package model

import (
	"gorm.io/gorm"
)

type Options struct {
	gorm.Model
	Value       string `json:"value" gorm:"type:text"`
	Question_id uint   `json:"question_id"`
	Is_right bool `json:"is_right"`
}