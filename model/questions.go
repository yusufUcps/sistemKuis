package model

import (
	"time"

	"gorm.io/gorm"
)

type Questions struct {
	gorm.Model
	Quiz_id uint `json:"quiz_id" gorm:"type:varchar(255)"`
	Question    string `json:"question" gorm:"type:text"`
	Options []Options `json:"options" gorm:"foreignKey:Question_id"`
}

type QuestionsRes struct{
	Id 				uint `json:"id"`
	Created_at 		time.Time`json:"created_at" `
    Updatad_at 		time.Time`json:"updated_at" `
	Quiz_id uint `json:"quiz_id" `
	Question    string `json:"question" `
	Options []Options `json:"options"`
}

type QuestionsResForQuiz struct{
	Id 				uint `json:"id"`
	Quiz_id			uint `json:"quiz_id" `
	Question    	string `json:"question" `
	Options []OptionsResForQuestion	`json:"options"`
}

func ConvertQuestionsRes(question *Questions) *QuestionsRes {
    questionRes := QuestionsRes{
        Id:         question.ID,
        Created_at: question.CreatedAt,
        Updatad_at: question.UpdatedAt,
		Question  : question.Question,
    }
    return &questionRes
}

func ConvertAllQuestionsQuiz(questions []Questions) []QuestionsResForQuiz {
	var questionRes []QuestionsResForQuiz

	for _, q := range questions {
		var optionRes []OptionsResForQuestion

		for _, opt := range q.Options {
			optionRes = append(optionRes, OptionsResForQuestion{
				Id:       opt.ID,
				Value:    opt.Value,
				Is_right: opt.Is_right,
			})
		}

		question := QuestionsResForQuiz{
			Id:      q.ID,
			Quiz_id: q.Quiz_id,
			Question: q.Question,
			Options:  optionRes,
		}

		questionRes = append(questionRes, question)
	}

	return questionRes
}

func GetQuestionIDs(questions []Questions) []uint {
	ids := make([]uint, len(questions))
	for i, question := range questions {
		ids[i] = question.ID
	}
	return ids
}

