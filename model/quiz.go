package model

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	User_id 			uint `json:"user_id"`
	Title     			string `json:"title" gorm:"type:varchar(225)"`
	Description    		string `json:"description" gorm:"type:text"`
	Start_date 			time.Time `json:"start_date" gorm:"type:datetime"`
	End_date 			time.Time `json:"end_date" gorm:"type:datetime"`
	Questions 			[]Questions `json:"questions" gorm:"foreignKey:Quiz_id"`
}

type QuizRes struct {
	Id 				uint `json:"id"`
	Created_at 		time.Time`json:"created_at" `
    Updatad_at 		time.Time`json:"updated_at" `
	User_id 		uint `json:"user_id"`
	Title    		string `json:"title" `
	Description    	string `json:"description" `
	Start_date 		time.Time `json:"start_date" gorm:"type:datetime"`
	End_date 		time.Time `json:"end_date" gorm:"type:datetime"`
}

func ConvertQuizRes(quiz *Quiz) *QuizRes {

    quizRes := QuizRes{
        Id:         quiz.ID,
        Created_at: quiz.CreatedAt,
        Updatad_at: quiz.UpdatedAt,
		User_id: quiz.User_id,
        Title : quiz.Title,
		Description: quiz.Description,
		Start_date: quiz.Start_date,
		End_date: quiz.End_date,
    }
    return & quizRes
}

func ConvertAllQuiz (quiz []Quiz) []QuizRes{
	
	var quizRes []QuizRes
	
	for i := range quiz {
		quizRes = append(quizRes, *ConvertQuizRes(&quiz[i]))
    }

	return quizRes
}