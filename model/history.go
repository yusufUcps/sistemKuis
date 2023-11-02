package model

import (
	"time"

	"gorm.io/gorm"
)

type Answers struct{
	Question_id uint `json:"question_id"`
	Option_id uint `json:"option_id"`
}

type HistoryScore struct {
	gorm.Model
	User_id uint `json:"user_id"`
	Name string `json:"name"` 
	Quiz_id uint   `json:"quiz_id"`
	Title string   `json:"title"`
	Right_answer int `json:"right_answer"`
	Wrong_answer int `json:"wrong_answer"`
	Score float64 `json:"score"`
}

type HistoryScoreRes struct {
	Id uint `json:"id"`
	User_id uint `json:"user_id"`
	Name string `json:"name"` 
	Quiz_id uint   `json:"quiz_id"`
	Title string   `json:"title"`
	Right_answer int `json:"right_answer"`
	Wrong_answer int `json:"wrong_answer"`
	Score float64 `json:"score"`
	Finish_at time.Time `json:"finish_at"`
}

type MyHistoryScoreRes struct {
	Id uint `json:"id"`
	Quiz_id uint   `json:"quiz_id"`
	Title string   `json:"title"`
	Right_answer int `json:"right_answer"`
	Wrong_answer int `json:"wrong_answer"`
	Score float64 `json:"score"`
	Finish_at time.Time `json:"finish_at"`
}

type HistoryScoreMyQuizRes struct {
	Id uint `json:"id"`
	User_id uint `json:"user_id"`
	Name string `json:"name"` 
	Right_answer int `json:"right_answer"`
	Wrong_answer int `json:"wrong_answer"`
	Score float64 `json:"score"`
	Finish_at time.Time `json:"finish_at"`
}

type HistoryAnswers struct {
	gorm.Model
	History_id uint `json:"history_id"`
	Question_id uint   `json:"qestion_id"`
	Question string `json:"question"`
	Option_id uint `json:"option_id"`
	Answer string `json:"answer"`
	Is_right bool `json:"is_right"`
}

type HistoryAnswersRes struct {
	Id uint `json:"id"`
	Question_id uint   `json:"qestion_id"`
	Question string `json:"question"`
	Option_id uint `json:"option_id"`
	Answer string `json:"answer"`
	Is_right bool `json:"is_right"`
}

type ScoreRes struct{
	Right_answer int `json:"right_answer"`
	Wrong_answer int `json:"wrong_answer"`
	Score float64 `json:"score"`
}



func GetOptionNQuestionIds(answers []Answers) ([]uint, []uint) {
	questionId := make([]uint, len(answers))
	optionId := make([]uint, len(answers))
	for i, answ := range answers {
		questionId[i] = answ.Question_id
		optionId[i] = answ.Option_id
	}
	return questionId, optionId
}

func GetQuizNUserIds(history []HistoryScore) ([]uint , []uint) {
	quizId := make([]uint, len(history))
	userId := make([]uint, len(history))
	for i, id := range history {
		quizId[i] = id.Quiz_id
		userId[i] = id.User_id
	}
	return userId, quizId
}

func ConvertAllMyHitoryScoreRes(histories []HistoryScore) []MyHistoryScoreRes {
    var myHistoryScoreRes []MyHistoryScoreRes

    for _, history := range histories {
        historyRes := MyHistoryScoreRes{
            Id:           history.ID,
            Quiz_id:      history.Quiz_id,
            Title:        history.Title,
            Right_answer: history.Right_answer,
            Wrong_answer: history.Wrong_answer,
            Score:        history.Score,
            Finish_at:    history.CreatedAt,
        }

        myHistoryScoreRes = append(myHistoryScoreRes, historyRes)
    }

    return myHistoryScoreRes
}

func ConvertAllHistoryScoreMyQuizRes(histories []HistoryScore) []HistoryScoreMyQuizRes {
    var historyScoreMyQuizRes []HistoryScoreMyQuizRes

    for _, history := range histories {
        historyRes := HistoryScoreMyQuizRes{
            Id:           history.ID,
            Right_answer: history.Right_answer,
            Wrong_answer: history.Wrong_answer,
            Score:        history.Score,
            Finish_at:    history.CreatedAt,
        }

        historyScoreMyQuizRes = append(historyScoreMyQuizRes, historyRes)
    }

    return historyScoreMyQuizRes
}


func GetScore(options []Options, questions []Questions, historyId uint) (ScoreRes, []HistoryAnswers) {
	var right_answer int = 0
	var score float64 = 0
	var wrong_answer int = 0

	var historyAnswer []HistoryAnswers 

	for i := range options {
		if options[i].Is_right {
			right_answer = right_answer + 1
		} else {
			wrong_answer = wrong_answer + 1
		}

		historyAnswer = append(historyAnswer, HistoryAnswers{
			History_id:  historyId,
			Question_id: options[i].Question_id,
			Option_id:   options[i].ID,
			Is_right:    options[i].Is_right,
			Question:    questions[i].Question,
			Answer:      options[i].Value,
		})
	}

	if right_answer+wrong_answer != 0 {
		score = float64(right_answer) / float64(wrong_answer+right_answer) * 100
	}

	var scoreRes ScoreRes

	scoreRes.Right_answer = right_answer
	scoreRes.Wrong_answer = wrong_answer
	scoreRes.Score = score

	return scoreRes, historyAnswer
}


func ConvertHistoryScore(score *HistoryScore) *HistoryScoreRes {
	var historyScoreRes HistoryScoreRes

	historyScoreRes.Id = score.ID
	historyScoreRes.User_id = score.User_id
	historyScoreRes.Quiz_id = score.Quiz_id
	historyScoreRes.Finish_at = score.CreatedAt 
	historyScoreRes.Name = score.Name
	historyScoreRes.Title = score.Title
	historyScoreRes.Right_answer = score.Right_answer
	historyScoreRes.Wrong_answer = score.Wrong_answer
	historyScoreRes.Score = score.Score

	return &historyScoreRes
}

func ConvertHistoryAnswer(answers []HistoryAnswers) []HistoryAnswersRes {
	var historyAnswerRes []HistoryAnswersRes

	for _, answer := range answers {
		var singleRes HistoryAnswersRes
		singleRes.Id = answer.ID
		singleRes.Question_id = answer.Question_id
		singleRes.Question = answer.Question
		singleRes.Answer = answer.Answer
		singleRes.Option_id = answer.Option_id
		singleRes.Is_right = answer.Is_right

		historyAnswerRes = append(historyAnswerRes, singleRes)
	}

	return historyAnswerRes
}
