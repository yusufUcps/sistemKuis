package repository

import (
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OptionsInterface interface {
	InsertOption(newOptions model.Options) (*model.Options, int)
	GetAllOptionsFromQuiz(questionsId uint) ([]model.Options, int)
	GetOptionByID(id uint) (*model.Options, int)
	UpdateOption(updateOptions model.Options, userId uint) (*model.Options, int)
	DeleteOption(questionsId uint, userId uint)  int
}

type OptionsModel struct {
	db *gorm.DB
}

func (qm *OptionsModel) Init(db *gorm.DB) {
	qm.db = db
}

func NewOptionsModel(db *gorm.DB) OptionsInterface {
	return &OptionsModel{
		db: db,
	}
}

func (om *OptionsModel) InsertOption(newOptions model.Options) (*model.Options, int) {

	if err := om.db.Create(&newOptions).Error; err != nil {
		logrus.Error("Repository: Insert data Option error, ", err.Error())
		return nil, 1
	}

	return &newOptions  , 0
}

func (om *OptionsModel) GetAllOptionsFromQuiz(questionsId uint) ([]model.Options, int) {
	var listOptions []model.Options
	var count int64

	if err := om.db.Where("question_id = ?", questionsId).Find(&listOptions).Count(&count).Error; err != nil {
		logrus.Error("Repository: Select method Get Options data error, ", err.Error())
		return nil, 1
	}

	if count == 0 {
		logrus.Error("Repository: not found Options")
		return nil, 2
	}

	return listOptions, 0
}

func (om *OptionsModel) GetOptionByID(id uint) (*model.Options, int) {

	var options model.Options

	if err := om.db.First(&options, id).Error; err != nil {
		logrus.Error("Repository: Get data options error, ", err.Error())
		return nil, 1
	}

	return &options, 0
}

func (om *OptionsModel) UpdateOption(updateOptions model.Options, userId uint) (*model.Options, int) {
	var options = model.Options{}
	var questions = model.Questions{}
	var quiz = model.Quiz{}

	if err := om.db.First(&options, updateOptions.ID).Error; err != nil {
		logrus.Error("Repository: Select method updateOptions data error, ", err.Error())
		return nil, 1
	}
	
	if err := om.db.First(&questions, options.Question_id).Error; err != nil {
		logrus.Error("Repository: Select method updateOptions data error, ", err.Error())
		return nil, 1
	}

	if err := om.db.First(&quiz, questions.Quiz_id).Error; err != nil {
		logrus.Error("Repository: Select method quiz in updateOptions data error, ", err.Error())
		return nil, 1
	}

	if userId != quiz.User_id {
		logrus.Error("Repository: UpdateOptions, Unauthorized")
		return nil, 2

	}

	options.Question_id = updateOptions.Question_id
	options.Value = updateOptions.Value
	options.Is_right = updateOptions.Is_right

	var qry = om.db.Save(&questions)
	if err := qry.Error; err != nil {
		logrus.Error("Repository: Save method updateOptions data error, ", err.Error())
		return nil, 1
	}

	return &options, 0
}

func (om *OptionsModel) DeleteOption(questionsId uint, userId uint)  int {

	var options = model.Options{}
	var questions = model.Questions{}
	var quiz = model.Quiz{}

	if err := om.db.First(&options, questionsId).Error; err != nil {
		logrus.Error("Repository: Select method updateOptions data error, ", err.Error())
		return 1
	}
	
	if err := om.db.First(&questions, options.Question_id).Error; err != nil {
		logrus.Error("Repository: Select method updateOptions data error, ", err.Error())
		return 1
	}

	if err := om.db.First(&quiz, questions.Quiz_id).Error; err != nil {
		logrus.Error("Repository: Select method quiz in updateOptions data error, ", err.Error())
		return 1
	}

	if userId != quiz.User_id {
		logrus.Error("Repository: UpdateOptions, Unauthorized")
		return 2

	}

	if err := om.db.Delete(&quiz).Error; err != nil {
		logrus.Error("Repository: Delete method DeleteOptiondata error, ", err.Error())
		return  1
	}

	return  0
}