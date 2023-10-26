package routes

import (
	"quiz/configs"
	"quiz/controller"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uc controller.UserController, cfg configs.ProgramConfig) {
	e.POST("/users", uc.Register())
	e.POST("/login", uc.Login())
	e.GET("/myprofile", uc.MyProfile(), echojwt.JWT([]byte(cfg.Secret)))
	e.PUT("/myprofile", uc.UpdateMyProfile(), echojwt.JWT([]byte(cfg.Secret)))
}