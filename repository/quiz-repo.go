package repository

import (
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QuizInterface interface {
	InsertQuiz(newQuiz model.Quiz) (*model.Quiz, int)
}

type QuizModel struct {
	db *gorm.DB
}

func (qm *QuizModel) Init(db *gorm.DB) {
	qm.db = db
}

func NewQuizModel(db *gorm.DB) QuizInterface {
	return &QuizModel{
		db: db,
	}
}

func (qm *QuizModel) InsertQuiz(newQuiz model.Quiz) (*model.Quiz, int) {

	if err := qm.db.Create(&newQuiz).Error; err != nil {
		logrus.Error("Repository: Insert data Quiz error, ", err.Error())
		return nil, 1
	}

	return &newQuiz, 0
}