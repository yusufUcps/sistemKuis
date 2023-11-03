package routes

import (
	"quiz/configs"
	"quiz/controller"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uc controller.UserControllInterface, cfg configs.ProgramConfig) {
	e.POST("/users", uc.Register())
	e.POST("/users/login", uc.Login())
	e.GET("/user", uc.MyProfile(), echojwt.JWT([]byte(cfg.Secret)))
	e.PUT("/user", uc.UpdateMyProfile(), echojwt.JWT([]byte(cfg.Secret)))
	e.DELETE("/user", uc.DeleteUser(), echojwt.JWT([]byte(cfg.Secret)))
}

func RouteQuiz(e *echo.Echo, uq controller.QuizControllInterface, cfg configs.ProgramConfig) {
	e.POST("/quizzes", uq.InsertQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/quizzes", uq.GetAllQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/quiz/:id", uq.GetQuizByID(), echojwt.JWT([]byte(cfg.Secret)))
	e.PUT("/quiz/:id", uq.UpdateQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.DELETE("/quiz/:id", uq.DeleteQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/quiz/me", uq.GetAllMyQuiz(), echojwt.JWT([]byte(cfg.Secret)))
}

func RouteQuestion(e *echo.Echo, uq controller.QuestionsControllInterface, cfg configs.ProgramConfig) { 
	e.POST("/question", uq.InsertQuestion(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/questions", uq.GetAllQuestionsQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/question/:id", uq.GetQuetionByID(), echojwt.JWT([]byte(cfg.Secret)))
	e.PUT("/question/:id", uq.UpdateQuestion(), echojwt.JWT([]byte(cfg.Secret)))
	e.POST("/questions/generate", uq.GenerateQuestion(), echojwt.JWT([]byte(cfg.Secret)))
	e.DELETE("/question/:id", uq.DeleteQuestion(), echojwt.JWT([]byte(cfg.Secret)))
}

func RouteOption(e *echo.Echo, uq controller.OptionsControllInterface, cfg configs.ProgramConfig) { 
	e.POST("/option", uq.InsertOption(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/options", uq.GetAllOptionsQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/option/:id", uq.GetOptionByID(), echojwt.JWT([]byte(cfg.Secret)))
	e.PUT("/option/:id", uq.UpdateOption(), echojwt.JWT([]byte(cfg.Secret)))
	e.DELETE("/option/:id", uq.DeleteOption(), echojwt.JWT([]byte(cfg.Secret)))
}

func RouteHistory(e *echo.Echo, uh controller.HistoryControllInterface, cfg configs.ProgramConfig) { 
	e.POST("/history/answers", uh.Answering(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/history/score", uh.GetAllMyHistoryScore(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/history/score/quiz", uh.GetAllHistoryScoreMyQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/history/score/:id", uh.GetHistoryScoreById(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/history/answer", uh.GetAllHistoryAnswer(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/history/score/export", uh.ExportMyHistoryScore(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/history/score/quiz/export", uh.ExportHistoryScoreMyQuiz(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/history/answer/export", uh.ExportHistoryAnswer(), echojwt.JWT([]byte(cfg.Secret)))
}
