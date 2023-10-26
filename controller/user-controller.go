package controller

import (
	"net/http"
	"quiz/configs"
	"quiz/helper"
	"quiz/model"
	"quiz/repository"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	cfg   configs.ProgramConfig
	model repository.UsersModel
}

func (uc *UserController) InitUserController(um repository.UsersModel, c configs.ProgramConfig) {
	uc.model = um
	uc.cfg = c
}

func (uc *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input model.Users
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil))
		}

		hashedPassword := helper.HashPassword(input.Password)
		input.Password = hashedPassword
		res, err := uc.model.Register(input)

		if err == 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Email already registered", nil))
		}

		if err == 2 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("cannot process data, something happend", nil))
		}

		var jwtToken = helper.GenerateJWT(uc.cfg.Secret, res.ID, res.Name)

		return c.JSON(http.StatusOK, helper.FormatResponseJWT("Succes create account", res, jwtToken))
	}
}
