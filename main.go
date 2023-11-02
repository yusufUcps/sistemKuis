package main

import (
	"fmt"
	"quiz/configs"
	"quiz/controller"
	"quiz/database"
	"quiz/helper"
	"quiz/repository"
	"quiz/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	var config = configs.InitConfig()

	db := database.InitModel(*config)
	database.Migrate(db)

	jwtInterface := helper.New(config.Secret)
	openAiInterface := helper.NewOpenAi(config.OpenAiKey)

	userModel := repository.NewUsersModel(db)
	quizModel := repository.NewQuizModel(db)
	questionModel := repository.NewQuestionsModel(db)

	userControll := controller.NewUserControllInterface(userModel, jwtInterface)
	quizControll := controller.NewQuizControllInterface(quizModel, jwtInterface)
	questionControll := controller.NewQuestionsControllInterface(questionModel, jwtInterface, openAiInterface)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
		}))

	routes.RouteUser(e, userControll, *config)
	routes.RouteQuiz(e, quizControll, *config)
	routes.RouteQuestion(e, questionControll, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}