package controller

import (
	"net/http"
	"quiz/helper"
	"quiz/model"
	"quiz/repository"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type QuestionsControllInterface interface {
	InsertQuestion() echo.HandlerFunc
	GetAllQuestionsQuiz() echo.HandlerFunc
	UpdateQuestion() echo.HandlerFunc
	GetQuetionByID() echo.HandlerFunc
	GenerateQuestion() echo.HandlerFunc
	DeleteQuestion() echo.HandlerFunc
}

type QuestionsController struct {
	repository repository.QuestionsInterface
	jwt helper.JWTInterface
	openAi helper.OpenAiInterface
}


func NewQuestionsControllInterface(r repository.QuestionsInterface, j helper.JWTInterface, o helper.OpenAiInterface) QuestionsControllInterface {
	return &QuestionsController{
		repository: r,
		jwt : j,
		openAi: o,
	}
}

func (qc *QuestionsController) InsertQuestion() echo.HandlerFunc {
	return func(c echo.Context) error {
		quizId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid quizId", nil, nil))
		}
		var newQuestions model.Questions

		newQuestions.Quiz_id = uint(quizId)

		if err := c.Bind(&newQuestions); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}
		
		res, errCase := qc.repository.InsertQuestion(newQuestions)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Insert Qusetion", nil, nil))
		}

		resConvert := model.ConvertQuestionsRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes create Question", resConvert, nil))
	}
}

func (qc *QuestionsController) GetAllQuestionsQuiz() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		quizId, err := strconv.Atoi(c.QueryParam("quizId"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid quizId", nil, nil))
		}

		page, err := strconv.Atoi(c.QueryParam("page"))
			if err != nil {
    		page = 1
		}

		pageSize, err := strconv.Atoi(c.QueryParam("pageSize")) 
			if err != nil {
    		pageSize = 10
		}

		quiz_id := uint(quizId) 
		
		res, count, errCase := qc.repository.GetAllQuestionsFromQuiz(page, pageSize, quiz_id)

		if errCase == 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid page", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All Questions Quiz", nil, nil))
		}

		if errCase == 3 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Question not Found", nil, nil))
		}

		resConvert := model.ConvertAllQuestionsQuiz(res)

		resPaging := model.ConvertPaging(page, pageSize, count) 

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get All Questions Quiz", resConvert, resPaging))
	}
}

func (qc *QuestionsController) GetQuetionByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		questionId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid questionId", nil, nil))
		}

		res, errCase := qc.repository.GetQuestionByID(uint(questionId))

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get Question Quiz", nil, nil))
		}

		resConvert := model.ConvertQuestionsRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get Question Quiz", resConvert, nil))
	}
}

func (qc *QuestionsController) UpdateQuestion() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")
		questionId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid questionId", nil, nil))
		}

		id := qc.jwt.ExtractToken(token.(*jwt.Token))

		var input = model.Questions{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}

		input.ID= uint(questionId)

		res, errCase := qc.repository.UpdateQuestion(input, id)


		if errCase == 2 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot update this question", nil, nil))
		}

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to update question", nil, nil))
		}

		resConvert := model.ConvertQuestionsRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("success update question", resConvert, nil))
	}
}

func (qc *QuestionsController) DeleteQuestion() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")
		questionId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid questionId", nil, nil))
		}

		id := qc.jwt.ExtractToken(token.(*jwt.Token))

		question_id := uint(questionId)

		errCase := qc.repository.DeleteQuestion(question_id, id)

		if errCase == 2 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot Delete this question", nil, nil))
		}

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to Delete question", nil, nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("success Delete question", nil, nil))
	}
}

func (qc *QuestionsController) GenerateQuestion() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		var newPromt model.OpenAiReq

		if err := c.Bind(&newPromt); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}
		
		newQuestions := qc.openAi.GenerateQuestions(newPromt)

		if newQuestions == nil{
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Generate Qusetions", nil, nil))
		}

		res, errCase := qc.repository.InsertGenerateQuestion(newQuestions)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Insert Generate Question", nil, nil))
		}

		resConvert := model.ConvertAllQuestionsQuiz(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Generate Question", resConvert, nil))
	}
}