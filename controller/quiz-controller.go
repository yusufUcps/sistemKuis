package controller

import (
	"net/http"
	"quiz/helper"
	"quiz/model"
	"quiz/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type QuizControllInterface interface {
	InsertQuiz() echo.HandlerFunc
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