package routes

import (
	"quiz/configs"
	"quiz/controller"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uc controller.UserControllInterface, cfg configs.ProgramConfig) {
	e.POST("/user", uc.Register())
	e.POST("/login", uc.Login())
	e.GET("/my-profile", uc.MyProfile(), echojwt.JWT([]byte(cfg.Secret)))
	e.PUT("/my-profile", uc.UpdateMyProfile(), echojwt.JWT([]byte(cfg.Secret)))
}

func RouteQuiz(e *echo.Echo, uq controller.QuizControllInterface, cfg configs.ProgramConfig) {
	e.POST("/quiz", uq.InsertQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/quiz", uq.GetAllQuiz(), echojwt.JWT([]byte(cfg.Secret)))
}