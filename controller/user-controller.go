package controller

import (
	"net/http"
	"quiz/configs"
	"quiz/helper"
	"quiz/model"
	"quiz/repository"

	"github.com/golang-jwt/jwt/v5"
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

		var jwtToken = helper.GenerateJWT(uc.cfg.Secret, res.ID)

		return c.JSON(http.StatusOK, helper.FormatResponseJWT("Succes create account", res, jwtToken))
	}
}

func (uc *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = model.Login{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil))
		}

		var res,err = uc.model.Login(input.Email, input.Password)

		if err == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("cannot process data, something happend", nil))
		}

		if err == 2 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("wrong email or password", nil))
		}

		var jwtToken = helper.GenerateJWT(uc.cfg.Secret, res.ID)

		if jwtToken == "" {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("cannot process data", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponseJWT("success", res ,jwtToken))
	}
}

func (uc *UserController) MyProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
	var token = c.Get("user")

	id := helper.ExtractToken(token.(*jwt.Token))

	res, err := uc.model.MyProfile(id)

	if err == 1 {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user profile")
	}
		return c.JSON(http.StatusOK, helper.FormatResponse("success get my profil", res))
	}
}

func (uc *UserController) UpdateMyProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")

		id := helper.ExtractToken(token.(*jwt.Token))

		var input = model.Users{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil))
		}

		hashedPassword := helper.HashPassword(input.Password)
		input.Password = hashedPassword
		input.ID = id

		res, err := uc.model.UpdateMyProfile(&input)

		if err == 1 {
			return echo.NewHTTPError(http.StatusNotFound, helper.FormatResponse("user profile not found", nil))
		}

		if err == 2 {
			return echo.NewHTTPError(http.StatusInternalServerError, helper.FormatResponse("failed to update user profile", nil))
		}

		if err == 3 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Email already registered", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("success update user profile", res))
	}
}


