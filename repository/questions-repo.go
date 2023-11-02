package repository

import (
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QuestionsInterface interface {
	InsertQuestion(newQuestion model.Questions) (*model.Questions, int)
	GetAllQuestionsFromQuiz(page int, pageSize int, quizId uint) ([]model.Questions, int64, int)
	GetQuestionByID(id uint) (*model.Questions, int)
	UpdateQuestion(updateQuestions model.Questions, userId uint) (*model.Questions, int)
	InsertGenerateQuestion(newQuestions []model.Questions) ([]model.Questions, int)
	DeleteQuestion(questionsId uint, userId uint)  int
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

func (qm *QuestionsModel) GetAllQuestionsFromQuiz(page int, pageSize int, quizId uint) ([]model.Questions, int64, int) {
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

func (qm *QuestionsModel) GetQuestionByID(id uint) (*model.Questions, int) {

	var questions model.Questions

	if err := qm.db.Preload("Options").First(&questions, id).Error; err != nil {
		logrus.Error("Repository: Get data Questions error, ", err.Error())
		return nil, 1
	}

	return &questions, 0
}

func (qm *QuestionsModel) UpdateQuestion(updateQuestions model.Questions, userId uint) (*model.Questions, int) {
	var questions = model.Questions{}
	var quiz = model.Quiz{}

	if err := qm.db.First(&questions, updateQuestions.ID).Error; err != nil {
		logrus.Error("Repository: Select method updateQuestions data error, ", err.Error())
		return nil, 1
	}

	if err := qm.db.First(&quiz, questions.Quiz_id).Error; err != nil {
		logrus.Error("Repository: Select method quiz in updateQuestions data error, ", err.Error())
		return nil, 1
	}

	if userId != quiz.User_id {
		logrus.Error("Repository: UpdateQuiz, Unauthorized")
		return nil, 2

	}

	questions.Quiz_id = updateQuestions.Quiz_id
	questions.Question = updateQuestions.Question
	questions.Options  = updateQuestions.Options

	var qry = qm.db.Save(&questions)
	if err := qry.Error; err != nil {
		logrus.Error("Repository: Save method updateQuestions data error, ", err.Error())
		return nil, 1
	}

	return &questions, 0
}

func (qm *QuestionsModel) InsertGenerateQuestion(newQuestions []model.Questions) ([]model.Questions, int) {

	if err := qm.db.Create(&newQuestions).Error; err != nil {
		logrus.Error("Repository: Insert data Generate Questions error, ", err.Error())
		return nil, 1
	}

	var questions []model.Questions

	if err := qm.db.Preload("Options").Where("id IN (?)", model.GetQuestionIDs(newQuestions)).Find(&questions).Error; err != nil {
		logrus.Error("Repository: Get data Generate Questions error, ", err.Error())
		return nil, 1
	}

	return questions , 0
}

func (qm *QuestionsModel) DeleteQuestion(questionsId uint, userId uint)  int {
	var questions = model.Questions{}
	var quiz = model.Quiz{}

	if err := qm.db.First(&questions, questionsId).Error; err != nil {
		logrus.Error("Repository: Select method quesrions in DeleteQuestions data error, ", err.Error())
		return  1
	}

	if err := qm.db.First(&quiz, questions.Quiz_id).Error; err != nil {
		logrus.Error("Repository: Select method quesrions in DeleteQuestions data error,", err.Error())
		return  1
	}

	if userId != quiz.User_id {
		logrus.Error("Repository: DeleteQuiz, Unauthorized")
		return  2
	}

	if err := qm.db.Delete(&questions).Error; err != nil {
		logrus.Error("Repository: Delete method DeleteQuiz data error, ", err.Error())
		return  1
	}

	return  0
}