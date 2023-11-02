package repository

import (
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QuestionsInterface interface {
	InsertQuestion(newQuestion model.Questions) (*model.Questions, int)
	GetAllQuestionsFromQuiz(page int, pageSize int, quizId uint) ([]model.Questions, int64, int)
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

func (qm *QuestionsModel) (page int, pageSize int, quizId uint) ([]model.Questions, int64, int) {
	var listQuestions []model.Questions
	var count int64

	if page <= 0 {
        return nil, 0, 1
    }

	offset := (page - 1) * pageSize

	if err := qm.db.Offset(offset).Limit(pageSize).Where("quiz_id = ?", quizId).Preload("Options").Find(&listQuestions).Count(&count).Error; err != nil {
		
		logrus.Error("Repository: Select method Get Questions data error, ", err.Error())
		return nil, 0, 2
	}

	if count == 0 {
		logrus.Error("Repository: not found Questions")
		return nil, 0, 3
	}

	return listQuestions, count, 0
}