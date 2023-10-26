package routes

import (
	"quiz/configs"
	"quiz/controller"
	
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uc controller.UserController, cfg configs.ProgramConfig) {
	e.POST("/users", uc.Register())
}