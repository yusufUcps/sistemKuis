package model

import (
	"time"

	"gorm.io/gorm"
)

type Options struct {
	gorm.Model
	Value       string `json:"value" gorm:"type:text"`
	Question_id uint   `json:"question_id"`
	Is_right 	bool `json:"is_right"`
}

type OptionsRes struct {
	Id 				uint `json:"id"`
	Created_at 		time.Time`json:"created_at" `
    Updatad_at 		time.Time`json:"updated_at" `
	Value       	string `json:"value"`
	Question_id 	uint   `json:"question_id"`
	Is_right 		bool `json:"is_right"`
}

type OptionsResForQuestion struct {
	Id 				uint `json:"id"`
	Value       	string `json:"value"`
	Is_right 		bool `json:"is_right"`
}


func ConvertOptionsRes(option *Options) *OptionsRes {
    optionRes := OptionsRes{
        Id:         option.ID,
        Created_at: option.CreatedAt,
        Updatad_at: option.UpdatedAt,
		Question_id: option.Question_id,
		Value:		option.Value,
		Is_right: 	option.Is_right,
    }
    return &optionRes
}

func ConvertOptionsResForQuestion(option *Options) *OptionsResForQuestion {
    optionRes := OptionsResForQuestion{
        Id:         option.ID,
		Value:		option.Value,
		Is_right: 	option.Is_right,
    }
    return &optionRes
}

func ConvertAllOptions (options []Options) []OptionsResForQuestion{
	
	var optionRes []OptionsResForQuestion
	
	for i := range options {
		optionRes = append(optionRes, *ConvertOptionsResForQuestion(&options[i]))
    }

	return optionRes
}

