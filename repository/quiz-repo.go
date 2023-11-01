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
	GetQuizByID(id uint) (*model.Quiz, int)
	UpdateQuiz(updateQuiz model.Quiz, userId uint) (*model.Quiz, int)
	DeleteQuiz(quizId uint, userId uint)  int
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

func (qm *QuizModel) GetQuizByID(id uint) (*model.Quiz, int) {

	var quiz model.Quiz

	if err := qm.db.Preload("Questions").First(&quiz, id).Error; err != nil {
		logrus.Error("Repository: Get data Quiz error, ", err.Error())
		return nil, 1
	}

	return &quiz, 0
}

func (qm *QuizModel) UpdateQuiz(updateQuiz model.Quiz, userId uint) (*model.Quiz, int) {
	var quiz = model.Quiz{}

	if err := qm.db.First(&quiz, updateQuiz.ID).Error; err != nil {
		logrus.Error("Repository: Select method UpdateQuiz data error, ", err.Error())
		return nil, 1
	}

	if userId != quiz.User_id {
		logrus.Error("Repository: UpdateQuiz, Unauthorized")
		return nil, 2

	}

	quiz.Title = updateQuiz.Title
	quiz.Description = updateQuiz.Description
	quiz.Start_date  = updateQuiz.Start_date 
	quiz.End_date  = updateQuiz.End_date 

	var qry = qm.db.Save(&quiz)
	if err := qry.Error; err != nil {
		logrus.Error("Repository: Save method UpdateQuiz data error, ", err.Error())
		return nil, 1
	}

	return &quiz, 0
}

func (qm *QuizModel) DeleteQuiz(quizId uint, userId uint)  int {
	var quiz = model.Quiz{}

	if err := qm.db.First(&quiz, quizId).Error; err != nil {
		logrus.Error("Repository: Select method UpdateQuiz data error, ", err.Error())
		return  1
	}

	if userId != quiz.User_id {
		logrus.Error("Repository: DeleteQuiz, Unauthorized")
		return  2
	}

	if err := qm.db.Delete(&quiz).Error; err != nil {
		logrus.Error("Repository: Delete method DeleteQuiz data error, ", err.Error())
		return  1
	}

	return  0
}