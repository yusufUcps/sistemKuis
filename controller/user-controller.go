package controller

import (
	"net/http"
	"quiz/helper"
	"quiz/model"
	"quiz/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserControllInterface interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	MyProfile() echo.HandlerFunc
	UpdateMyProfile() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
}

type UserController struct {
	repository repository.UsersInterface
	jwt helper.JWTInterface
}


func NewUserControllInterface(r repository.UsersInterface, j helper.JWTInterface) UserControllInterface {
	return &UserController{
		repository: r,
		jwt: j,
	}
}

func (uc *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input model.Users
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}
		
		res, errCase := uc.repository.Register(input)

		if errCase == 1 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Email already registered", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed Register User", nil, nil))
		}

		var jwtToken = uc.jwt.GenerateJWT(res.ID)
		if jwtToken == "" {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("cannot process jwt token", nil, nil))
		}

		resConvert := model.ConvertRegisterRes(res, jwtToken)

		return c.JSON(http.StatusOK, helper.FormatResponse("Succes create account", resConvert, nil))
	}
}

func (uc *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = model.Login{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}

		var res, errCase = uc.repository.Login(input.Email, input.Password)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Login", nil, nil))
		}

		if errCase == 2 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("wrong email or password", nil, nil))
		}

		var jwtToken = uc.jwt.GenerateJWT(res.ID)

		if jwtToken == "" {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("cannot process jwt token", nil, nil))
		}

		resConvert := model.ConvertLoginRes(res, jwtToken)

		return c.JSON(http.StatusOK, helper.FormatResponse("Login Success", resConvert ,nil))
	}
}

func (uc *UserController) MyProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")

		id := uc.jwt.ExtractToken(token.(*jwt.Token))

		res, errCase := uc.repository.MyProfile(id)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to Get MyProfile", nil, nil))
		}

		resConvert := model.ConvertMyProfileRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("success get my profil", resConvert, nil))
	}
}

func (uc *UserController) UpdateMyProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")

		id := uc.jwt.ExtractToken(token.(*jwt.Token))

		var input = model.Users{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil, nil))
		}

		input.ID = id

		res, errCase := uc.repository.UpdateMyProfile(input)


		if errCase == 2 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Email already registered", nil, nil))
		}

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to update user profile", nil, nil))
		}

		resConvert := model.ConvertMyProfileRes(res)

		return c.JSON(http.StatusOK, helper.FormatResponse("success update user profile", resConvert, nil))
	}
}

func (uc *UserController) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var token = c.Get("user")

		id := uc.jwt.ExtractToken(token.(*jwt.Token))


		errCase := uc.repository.DeleteUser(id)

		if errCase == 1 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("failed to Delete User", nil, nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("success Delete User", nil, nil))
	}
}
