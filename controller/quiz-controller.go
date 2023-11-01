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

type QuizControllInterface interface {
	InsertQuiz() echo.HandlerFunc
	GetAllQuiz() echo.HandlerFunc
	GetQuizByID() echo.HandlerFunc
	UpdateQuiz() echo.HandlerFunc
}

type QuizController struct {
	repository repository.QuizInterface
	jwt helper.JWTInterface
}


func NewQuizControllInterface(r repository.QuizInterface, j helper.JWTInterface) QuizControllInterface {
	return &QuizController{
		repository: r,
		jwt : j,
	}
}

func (qc *QuizController) InsertQuiz() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")
		var newQuiz model.Quiz

		id := qc.jwt.ExtractToken(token.(*jwt.Token))

		newQuiz.User_id = id

		if err := c.Bind(&newQuiz); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}
		
		res, errCase := qc.repository.InsertQuiz(newQuiz)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Insert Quiz", nil, nil))
		}

		resConvert := model.ConvertQuizRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes create Quiz", resConvert, nil))
	}
}

func (qc *QuizController) GetAllQuiz() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		search := c.QueryParam("search")

		page, err := strconv.Atoi(c.QueryParam("page"))
			if err != nil {
    		page = 1
		}

		pageSize, err := strconv.Atoi(c.QueryParam("pageSize")) 
			if err != nil {
    		pageSize = 10
		}
		
		res, count, errCase := qc.repository.GetAllQuiz(page, pageSize, search)

		if errCase == 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid page", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All Quiz", nil, nil))
		}

		if errCase == 3 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Quiz not Found", nil, nil))
		}

		resConvert := model.ConvertAllQuiz(res)

		resPaging := model.ConvertPaging(page, pageSize, count) 

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get All Quiz", resConvert, resPaging))
	}
}

func (qc *QuizController) GetQuizByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		quizId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid quizId", nil, nil))
		}

		res, errCase := qc.repository.GetQuizByID(uint(quizId))

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get Quiz", nil, nil))
		}

		resConvert := model.ConvertQuizRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get Quiz", resConvert, nil))
	}
}

func (qc *QuizController) UpdateQuiz() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")
		quizId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid quizId", nil, nil))
		}

		id := qc.jwt.ExtractToken(token.(*jwt.Token))

		var input = model.Quiz{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}

		input.ID= uint(quizId)

		res, errCase := qc.repository.UpdateQuiz(input, id)


		if errCase == 2 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot update this quiz", nil, nil))
		}

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to update Quiz", nil, nil))
		}

		resConvert := model.ConvertQuizRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("success update Quiz", resConvert, nil))
	}
}