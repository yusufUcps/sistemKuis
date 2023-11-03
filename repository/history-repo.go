package repository

import (
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HistoryInterface interface {
	AnswersInsert(answers []model.Answers, user_id uint, quiz_id uint) (*model.HistoryScore, int)
	GetAllMyHistoryScore(page int, pageSize int, userId uint, search string) ([]model.HistoryScore, int64, int)
	GetAllHistoryScoreMyQuiz(page int, pageSize int, quizId uint, search string, userId uint) ([]model.HistoryScore, int64, int)
	GetHistoryScoreById(historyId uint, userId uint) (*model.HistoryScore, int)
	GetAllHistoryAnswer(page int, pageSize int, historyId uint, userId uint) ([]model.HistoryAnswers, int64, int)
	ExMyHistoryScore(userId uint) ([]model.HistoryScore, int)
	ExHistoryScoreMyQuiz(quizId uint, userId uint) ([]model.HistoryScore, int)
	ExHistoryAnswer(historyId uint, userId uint) ([]model.HistoryAnswers, int)
}

type HistoryModel struct {
	db *gorm.DB
}

func (qm *HistoryModel) Init(db *gorm.DB) {
	qm.db = db
}

func NewHistoryModel(db *gorm.DB) HistoryInterface {
	return &HistoryModel{
		db: db,
	}
}

func (hm *HistoryModel) AnswersInsert(answers []model.Answers, user_id uint, quiz_id uint) (*model.HistoryScore, int) {

	var historyScore model.HistoryScore

	historyScore.User_id = user_id
	historyScore.Quiz_id = quiz_id

	var quiz1 model.Quiz

	if err := hm.db.First(&quiz1, quiz_id).Error; err != nil {
		logrus.Error("Repository: Select method quiz in updateOptions data error, ", err.Error())
		return nil,1
	}

	questionId, optionId := model.GetOptionNQuestionIds(answers)

	var questions []model.Questions
	if err := hm.db.Where("id IN (?)", questionId).Find(&questions).Error; err != nil {
		logrus.Error("Repository: select method question in answer data error, ", err.Error())
		return nil, 1
	}

	var options []model.Options

	if err := hm.db.Where("id IN (?)", optionId).Find(&options).Error; err != nil {
		logrus.Error("Repository: select method option in answer data error, ", err.Error())
		return nil, 1
	}

	if questions[1].Quiz_id != quiz_id{
		logrus.Error("Repository: invalid quiz id, ")
		return nil, 1
	}


	if err := hm.db.Create(&historyScore).Error; err != nil {
		logrus.Error("Repository: Insert data historyScore error, ", err.Error())
		return nil, 1
	}


	scoreQuiz , historyAnswer := model.GetScore(options, questions, historyScore.ID)

	if err := hm.db.Create(&historyAnswer).Error; err != nil {
		logrus.Error("Repository: Insert data historyAnswer error, ", err.Error())
		return nil, 1
	}

	var quiz model.Quiz

	if err := hm.db.Select("title").Where("id = ?", historyScore.Quiz_id).Find(&quiz).Error; err != nil {
		logrus.Error("Repository: select method title in answer data error, ", err.Error())
		return nil, 1
	}

	var user model.Users

	if err := hm.db.Select("name").Where("id = ?", historyScore.User_id).Find(&user).Error; err != nil {
		logrus.Error("Repository: select method name in answer data error, ", err.Error())
		return nil, 1
	}

	historyScore.Wrong_answer = scoreQuiz.Wrong_answer
	historyScore.Right_answer = scoreQuiz.Right_answer
	historyScore.Score = scoreQuiz.Score 
	historyScore.Title = quiz.Title
	historyScore.Name = user.Name

	var qry = hm.db.Save(&historyScore)
	if err := qry.Error; err != nil {
		logrus.Error("Repository: Save method historyScore data error, ", err.Error())
		return nil, 1
	}
	
	return &historyScore , 0
}

func (hm *HistoryModel) GetAllMyHistoryScore(page int, pageSize int, userId uint, search string) ([]model.HistoryScore, int64, int) {
	var listHistoryScore []model.HistoryScore
	var count int64

	if page <= 0 {
        return nil, 0, 1
    }

	offset := (page - 1) * pageSize

	if err := hm.db.Offset(offset).Limit(pageSize).Where("user_id = ? AND title LIKE ?", userId, "%"+search+"%").Find(&listHistoryScore).Count(&count).Error; err != nil {
		
		logrus.Error("Repository: Select method Get MyHistoryScore data error, ", err.Error())
		return nil, 0, 2
	}

	if count == 0 {
		logrus.Error("Repository: not found MyHistoryScore")
		return nil, 0, 3
	}

	return listHistoryScore, count, 0
}

func (hm *HistoryModel) GetAllHistoryScoreMyQuiz(page int, pageSize int, quizId uint, search string, userId uint) ([]model.HistoryScore, int64, int) {
	var listHistoryScore []model.HistoryScore
	var count int64

	if page <= 0 {
        return nil, 0, 1
    }

	offset := (page - 1) * pageSize

	var quiz = model.Quiz{}
	if err := hm.db.First(&quiz, quizId).Error; err != nil {
		logrus.Error("Repository: Select method UpdateQuiz data error, ", err.Error())
		return nil, 0, 2
	}

	if userId != quiz.User_id{
		logrus.Error("Repository: GetAllHistoryScoreMyQuiz, Unauthorized")
		return  nil, 0, 4
	}

	if err := hm.db.Offset(offset).Limit(pageSize).Where("quiz_id = ? AND name LIKE ?", quizId, "%"+search+"%").Find(&listHistoryScore).Count(&count).Error; err != nil {
		
		logrus.Error("Repository: Select method Get listHistoryScoreMyQuiz data error, ", err.Error())
		return nil, 0, 2
	}

	if count == 0 {
		logrus.Error("Repository: not found HistoryScoreMyQuiz")
		return nil, 0, 3
	}

	return listHistoryScore, count, 0
}

func (hm *HistoryModel) GetHistoryScoreById(historyId uint, userId uint) (*model.HistoryScore, int) {
	var HistoryAnswers model.HistoryScore

	if err := hm.db.First(&HistoryAnswers, historyId).Error; err != nil {
		logrus.Error("Repository: Select method HistoryScore data error, ", err.Error())
		return nil, 1
	}

	var quiz = model.Quiz{}
	if err := hm.db.First(&quiz, HistoryAnswers.Quiz_id).Error; err != nil {
		logrus.Error("Repository: Select method quiz in HistoryScore data error, ", err.Error())
		return nil, 1
	}

	if userId != quiz.User_id && userId != HistoryAnswers.User_id{
		logrus.Error("Repository: HistoryScore, Unauthorized")
		return  nil, 2
	}

	return  &HistoryAnswers, 0
}

func (hm *HistoryModel) GetAllHistoryAnswer(page int, pageSize int, historyId uint, userId uint) ([]model.HistoryAnswers, int64, int) {
	var listHistoryAnswers []model.HistoryAnswers
	var count int64

	if page <= 0 {
        return nil, 0, 1
    }

	var history = model.HistoryScore{}
	if err := hm.db.First(&history, historyId).Error; err != nil {
		logrus.Error("Repository: Select method historyScore data error, ", err.Error())
		return nil, 0, 2
	}

	var quiz = model.Quiz{}
	if err := hm.db.First(&quiz, history.Quiz_id).Error; err != nil {
		logrus.Error("Repository: Select method Quiz in history data error, ", err.Error())
		return nil, 0, 2
	}

	if userId != quiz.User_id && userId != history.User_id{
		logrus.Error("Repository:  HistoryAnswer, Unauthorized")
		return  nil, 0, 4
	}

	offset := (page - 1) * pageSize

	if err := hm.db.Offset(offset).Limit(pageSize).Where("history_id = ?", historyId).Find(&listHistoryAnswers).Count(&count).Error; err != nil {
		
		logrus.Error("Repository: Select method Get HistoryAnswer data error, ", err.Error())
		return nil, 0, 2
	}

	if count == 0 {
		logrus.Error("Repository: not found HistoryAnswer")
		return nil, 0, 3
	}

	return  listHistoryAnswers, count, 0
}

func (hm *HistoryModel) ExMyHistoryScore(userId uint) ([]model.HistoryScore, int) {
	var listHistoryScore []model.HistoryScore
	var count int64

	if err := hm.db.Where("user_id = ?", userId).Find(&listHistoryScore).Count(&count).Error; err != nil {
		
		logrus.Error("Repository: Select method Get ExlistHistoryScore data error, ", err.Error())
		return nil, 1
	}

	if count == 0 {
		logrus.Error("Repository: not found Questions")
		return nil, 2
	}

	return listHistoryScore, 0
}

func (hm *HistoryModel) ExHistoryScoreMyQuiz(quizId uint, userId uint) ([]model.HistoryScore, int) {
	var listHistoryScore []model.HistoryScore
	var count int64

	var quiz = model.Quiz{}
	if err := hm.db.First(&quiz, quizId).Error; err != nil {
		logrus.Error("Repository: Select method UpdateQuiz data error, ", err.Error())
		return nil,  1
	}

	if userId != quiz.User_id{
		logrus.Error("Repository: GetAllHistoryScoreMyQuiz, Unauthorized")
		return  nil, 3
	}

	if err := hm.db.Where("quiz_id = ?", quizId).Find(&listHistoryScore).Count(&count).Error; err != nil {
		
		logrus.Error("Repository: Select method Get listHistoryScoreMyQuiz data error, ", err.Error())
		return nil, 1
	}

	if count == 0 {
		logrus.Error("Repository: not found Questions")
		return nil, 2
	}

	return listHistoryScore, 0
}

func (hm *HistoryModel) ExHistoryAnswer(historyId uint, userId uint) ([]model.HistoryAnswers, int) {
	var listHistoryAnswers []model.HistoryAnswers
	var count int64

	var history = model.HistoryScore{}
	if err := hm.db.First(&history, historyId).Error; err != nil {
		logrus.Error("Repository: Select method UpdateQuiz data error, ", err.Error())
		return nil,  1
	}

	var quiz = model.Quiz{}
	if err := hm.db.First(&quiz, history.Quiz_id).Error; err != nil {
		logrus.Error("Repository: Select method UpdateQuiz data error, ", err.Error())
		return nil,  1
	}

	if userId != quiz.User_id && userId != history.User_id{
		logrus.Error("Repository: DeleteQuiz, Unauthorized")
		return  nil,  2
	}


	if err := hm.db.Where("history_id = ?", historyId).Find(&listHistoryAnswers).Count(&count).Error; err != nil {
		
		logrus.Error("Repository: Select method Get listHistoryScoreMyQuiz data error, ", err.Error())
		return nil,  2
	}

	if count == 0 {
		logrus.Error("Repository: not found Questions")
		return nil, 3
	}

	return  listHistoryAnswers,  0
}