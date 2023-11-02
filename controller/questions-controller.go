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
	
}

type QuestionsController struct {
	repository repository.QuestionsInterface
	jwt helper.JWTInterface
}


func NewQuestionsControllInterface(r repository.QuestionsInterface, j helper.JWTInterface) QuestionsControllInterface {
	return &QuestionsController{
		repository: r,
		jwt : j,
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