package repository

import (
	"quiz/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QuizInterface interface {
	InsertQuiz(newQuiz model.Quiz) (*model.Quiz, int)
	GetAllQuiz(page int, pageSize int, search string) ([]model.Quiz, int64, int)
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

func (qm *QuizModel) GetAllQuiz(page int, pageSize int, search string) ([]model.Quiz, int64, int) {
	var listQuiz []model.Quiz
	var count int64

	if page <= 0 {
        return nil, 0, 1
    }

	offset := (page - 1) * pageSize

	if err := qm.db.Where("title LIKE ? AND start_date <= ? AND end_date >= ?", "%"+search+"%", time.Now(), time.Now()).Offset(offset).Limit(pageSize).Find(&listQuiz).Count(&count).Error; err != nil {
		logrus.Error("Repository: Select method UpdateQuiz data error, ", err.Error())
		return nil, 0, 1
	}

	if count == 0 {
		logrus.Error("Repository: not found Quiz")
		return nil, 0, 3
	}

	return listQuiz, count, 0
}