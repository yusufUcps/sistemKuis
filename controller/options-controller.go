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

