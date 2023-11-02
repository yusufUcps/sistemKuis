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

type HistoryControllInterface interface {
	Answering() echo.HandlerFunc
	GetAllMyHistoryScore() echo.HandlerFunc
	
}

type HistoryController struct {
	repository repository.HistoryInterface
	jwt helper.JWTInterface
}


func NewHistoryControllInterface(r repository.HistoryInterface, j helper.JWTInterface) HistoryControllInterface {
	return &HistoryController{
		repository: r,
		jwt : j,
	}
}

func (hc *HistoryController) Answering() echo.HandlerFunc {
	return func(c echo.Context) error {

		var token = c.Get("user")
		id := hc.jwt.ExtractToken(token.(*jwt.Token))

		quizId, err := strconv.Atoi(c.QueryParam("quizId"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid quizId", nil, nil))
		}

		var newAnswers []model.Answers

		if err := c.Bind(&newAnswers); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}
		
		res, errCase := hc.repository.AnswersInsert(newAnswers ,id , uint(quizId))

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Insert Qusetion", nil, nil))
		}

		resConvert := model.ConvertHistoryScore(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes save asswer", resConvert, nil))
	}
}

func (hc *HistoryController) GetAllMyHistoryScore() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		var token = c.Get("user")
		id := hc.jwt.ExtractToken(token.(*jwt.Token))

		search := c.QueryParam("search")

		page, err := strconv.Atoi(c.QueryParam("page"))
			if err != nil {
    		page = 1
		}

		pageSize, err := strconv.Atoi(c.QueryParam("pageSize")) 
			if err != nil {
    		pageSize = 10
		}
		
		res, count, errCase := hc.repository.GetAllMyHistoryScore(page, pageSize, id, search)

		if errCase == 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid page", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All My History Score", nil, nil))
		}

		if errCase == 3 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("History Score not Found", nil, nil))
		}

		resConvert := model.ConvertAllMyHitoryScoreRes(res)

		resPaging := model.ConvertPaging(page, pageSize, count) 

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get All My History Score", resConvert, resPaging))
	}
}