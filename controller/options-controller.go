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
	UpdateOption() echo.HandlerFunc
	DeleteOption() echo.HandlerFunc
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
		
		questionsId, err := strconv.Atoi(c.QueryParam("questionId"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid questionId", nil, nil))
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

func (op *OptionsController) UpdateOption() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")
		optionId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid optionId", nil, nil))
		}

		id := op.jwt.ExtractToken(token.(*jwt.Token))

		var input = model.Options{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}

		input.ID= uint(optionId)

		res, errCase := op.repository.UpdateOption(input, id)


		if errCase == 2 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot update this option", nil, nil))
		}

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to update option", nil, nil))
		}

		resConvert := model.ConvertOptionsRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("success update option", resConvert, nil))
	}
}

func (op *OptionsController) DeleteOption() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")
		optionId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid quizId", nil, nil))
		}

		id := op.jwt.ExtractToken(token.(*jwt.Token))

		option_id := uint(optionId)

		errCase := op.repository.DeleteOption(option_id, id)

		if errCase == 2 {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You cannot Delete this Option", nil, nil))
		}

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to Delete Option", nil, nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("success Delete Option", nil, nil))
	}
}
