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
	GetAllHistoryScoreMyQuiz() echo.HandlerFunc
	GetHistoryScoreById() echo.HandlerFunc
	GetAllHistoryAnswer() echo.HandlerFunc
	ExportMyHistoryScore() echo.HandlerFunc
	ExportHistoryScoreMyQuiz() echo.HandlerFunc
	ExportHistoryAnswer() echo.HandlerFunc
}

type HistoryController struct {
	repository repository.HistoryInterface
	jwt helper.JWTInterface
	Export helper.ExportInterface
}


func NewHistoryControllInterface(r repository.HistoryInterface, j helper.JWTInterface, e helper.ExportInterface) HistoryControllInterface {
	return &HistoryController{
		repository: r,
		jwt : j,
		Export: e,
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

		if errCase == 2 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid quizId", nil, nil))
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

func (hc *HistoryController) GetAllHistoryScoreMyQuiz() echo.HandlerFunc {
	return func(c echo.Context) error {
		search := c.QueryParam("search")
		var token = c.Get("user")
		id := hc.jwt.ExtractToken(token.(*jwt.Token))
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
		
		res, count, errCase := hc.repository.GetAllHistoryScoreMyQuiz(page, pageSize, uint(quizId), search, id)

		if errCase == 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid page", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All History Score MyQuiz", nil, nil))
		}

		if errCase == 3 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("HistoryScore not Found", nil, nil))
		}

		if errCase == 4 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot Get this History", nil, nil))
		}

		resConvert := model.ConvertAllHistoryScoreMyQuizRes(res)

		resPaging := model.ConvertPaging(page, pageSize, count) 

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get Get All History Score MyQuiz", resConvert, resPaging))
	}
}

func (hc *HistoryController) GetHistoryScoreById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")
		id := hc.jwt.ExtractToken(token.(*jwt.Token))
		
		historyId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid historyId", nil, nil))
		}

		res, errCase := hc.repository.GetHistoryScoreById(uint(historyId),id)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get history Score", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot Get this History", nil, nil))
		}

		resConvert := model.ConvertHistoryScore(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get history score", resConvert, nil))
	}
}

func (hc *HistoryController) GetAllHistoryAnswer() echo.HandlerFunc {
	return func(c echo.Context) error {

		var token = c.Get("user")
		id := hc.jwt.ExtractToken(token.(*jwt.Token))

		historyId, err := strconv.Atoi(c.QueryParam("historyId"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid historyId", nil, nil))
		}

		page, err := strconv.Atoi(c.QueryParam("page"))
			if err != nil {
    		page = 1
		}

		pageSize, err := strconv.Atoi(c.QueryParam("pageSize")) 
			if err != nil {
    		pageSize = 10
		}
		
		res, count, errCase := hc.repository.GetAllHistoryAnswer(page, pageSize, uint(historyId), id)

		if errCase == 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid page", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All History Answer", nil, nil))
		}

		if errCase == 3 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("HistoryAnswer not Found", nil, nil))
		}

		if errCase == 4 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot Get this History", nil, nil))
		}

		resConvert := model.ConvertHistoryAnswer(res)

		resPaging := model.ConvertPaging(page, pageSize, count) 

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get Get All History Answer", resConvert, resPaging))
	}
}

func (hc *HistoryController) ExportMyHistoryScore() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		var token = c.Get("user")
		id := hc.jwt.ExtractToken(token.(*jwt.Token))
		
		res, errCase := hc.repository.ExMyHistoryScore(id)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All My History Score", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("History Score not Found", nil, nil))
		}

		resConvert := model.ConvertAllMyHitoryScoreRes(res)

		resExport, err := hc.Export.ExportMyHistoryScore(resConvert, id)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to export", nil, nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Export My History Score", resExport, nil))
	}
}

func (hc *HistoryController) ExportHistoryScoreMyQuiz() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		var token = c.Get("user")
		id := hc.jwt.ExtractToken(token.(*jwt.Token))

		quizId, err := strconv.Atoi(c.QueryParam("quizId"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid quizId", nil, nil))
		}
		
		res, errCase := hc.repository.ExHistoryScoreMyQuiz(uint(quizId), id)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All My History Score", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("History Score not Found", nil, nil))
		}

		if errCase == 3 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot export this History", nil, nil))
		}

		resConvert := model.ConvertAllHistoryScoreMyQuizRes(res)

		resExport, err := hc.Export.ExportHistoryScoreMyQuiz(resConvert, uint(quizId))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to export", nil, nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Export History Score Quiz", resExport, nil))
	}
}

func (hc *HistoryController) ExportHistoryAnswer() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		var token = c.Get("user")
		id := hc.jwt.ExtractToken(token.(*jwt.Token))

		hitoriId, err := strconv.Atoi(c.QueryParam("historyId"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid historyId", nil, nil))
		}
		
		res, errCase := hc.repository.ExHistoryAnswer(uint(hitoriId), id)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All My History Answer", nil, nil))
		}

		if errCase == 3 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("History Answer not Found", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot export this History", nil, nil))
		}

		resConvert := model.ConvertHistoryAnswer(res)

		resExport, err := hc.Export.ExportHistoryAnswer(resConvert, id)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to export", nil, nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes export HistoryAnswer", resExport, nil))
	}
}