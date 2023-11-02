package repository

import (
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QuestionsInterface interface {
	InsertQuestion(newQuestion model.Questions) (*model.Questions, int)
}

type QuestionsModel struct {
	db *gorm.DB
}

func (qm *QuestionsModel) Init(db *gorm.DB) {
	qm.db = db
}

func NewQuestionsModel(db *gorm.DB) QuestionsInterface {
	return &QuestionsModel{
		db: db,
	}
}

func (qm *QuestionsModel) InsertQuestion(newQuestion model.Questions) (*model.Questions, int) {

	if err := qm.db.Create(&newQuestion).Error; err != nil {
		logrus.Error("Repository: Insert data Questions error, ", err.Error())
		return nil, 1
	}

	var question model.Questions

	if err := qm.db.Preload("Options").First(&question, newQuestion.ID).Error; err != nil {
		logrus.Error("Repository: Get data Questions error, ", err.Error())
		return nil, 1
	}

	return &question , 0
}