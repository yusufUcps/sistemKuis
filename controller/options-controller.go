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

type OptionsControllInterface interface {
	InsertOption() echo.HandlerFunc
	GetAllOptionsQuiz() echo.HandlerFunc
	GetOptionByID() echo.HandlerFunc
	
}

type OptionsController struct {
	repository repository.OptionsInterface
	jwt helper.JWTInterface
}


func NewOptionsControllInterface(r repository.OptionsInterface, j helper.JWTInterface) OptionsControllInterface {
	return &OptionsController{
		repository: r,
		jwt : j,
	}
}

func (op *OptionsController) InsertOption() echo.HandlerFunc {
	return func(c echo.Context) error {
		var newOptions model.Options

		if err := c.Bind(&newOptions); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}
		
		res, errCase := op.repository.InsertOption(newOptions)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Insert Qusetion", nil, nil))
		}

		resConvert := model.ConvertOptionsRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes create Question", resConvert, nil))
	}
}

func (op *OptionsController) GetAllOptionsQuiz() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		questionsId, err := strconv.Atoi(c.QueryParam("QuestionId"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid QuestionId", nil, nil))
		}

		options_id := uint(questionsId) 
		
		res, errCase := op.repository.GetAllOptionsFromQuiz(options_id)

		if errCase == 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid page", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get All Options Question", nil, nil))
		}

		if errCase == 3 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Options not Found", nil, nil))
		}

		resConvert := model.ConvertAllOptions(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get All Option", resConvert,nil))
	}
}

func (op *OptionsController) GetOptionByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		optionId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid optionId", nil, nil))
		}

		res, errCase := op.repository.GetOptionByID(uint(optionId))

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get Option Question", nil, nil))
		}

		resConvert := model.ConvertOptionsRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes Get Option Question", resConvert, nil))
	}
}