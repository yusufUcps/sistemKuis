package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quiz/mocks"
	"quiz/model"
	"strings"
	"time"

	"testing"

	"quiz/controller"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertQuiz(t *testing.T) {
	t.Run("SuccessfulInsertQuiz", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		userID := uint(1)

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID

		tokenString, err := token.SignedString([]byte("s3cr3t"))
		if err != nil {
			t.Errorf("token signing error: %s", err)
		}

		start := time.Date(2023, time.November, 3, 8, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 0, 7)

		newQuiz := model.Quiz{
			User_id:     userID,
			Title:       "Quiz Title",
			Description: "Description of the quiz",
			Start_date:  start,
			End_date:    end,
			Questions:   nil,
		}

		mockJWT.On("ExtractToken", mock.Anything).Return(userID)
		mockRepo.On("InsertQuiz", newQuiz).Return(&newQuiz, 0)

		e := echo.New()
		reqJSON, err := json.Marshal(newQuiz)
		if err != nil {
			t.Errorf("JSON marshaling error: %s", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/quizzes", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)

		err = quizController.InsertQuiz()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"Succes create Quiz\"")

		mockJWT.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
	t.Run("FailedInsertQuizInvalidInput", func(t *testing.T) {

		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		userID := uint(1)

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID

		tokenString, err := token.SignedString([]byte("s3cr3t"))
		if err != nil {
			t.Errorf("token signing error: %s", err)
		}

		mockJWT.On("ExtractToken", mock.Anything).Return(userID)
		invalidInput := `{"email": "test@example.com"`
		e := echo.New()

		if err != nil {
			t.Errorf("JSON marshaling error: %s", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/quizzes", strings.NewReader(invalidInput))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)

		err = quizController.InsertQuiz()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"invalid user input\"")

		mockJWT.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FailedInsertQuiz", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		userID := uint(1)

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID

		tokenString, err := token.SignedString([]byte("s3cr3t"))
		if err != nil {
			t.Errorf("token signing error: %s", err)
		}

		start := time.Date(2023, time.November, 3, 8, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 0, 7)

		newQuiz := model.Quiz{
			User_id:     userID,
			Title:       "Quiz Title",
			Description: "Description of the quiz",
			Start_date:  start,
			End_date:    end,
			Questions:   nil,
		}

		mockJWT.On("ExtractToken", mock.Anything).Return(userID)
		mockRepo.On("InsertQuiz", newQuiz).Return(nil, 1)

		e := echo.New()
		reqJSON, err := json.Marshal(newQuiz)
		if err != nil {
			t.Errorf("JSON marshaling error: %s", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/quizzes", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)

		err = quizController.InsertQuiz()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"Failed to Insert Quiz\"")

		mockJWT.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

}

func TestGetAllQuiz(t *testing.T) {
	t.Run("SuccessfulGetAllQuiz", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		fakeQuizData := []model.Quiz{
			{
				User_id:     1,
				Title:       "Sample Quiz 1",
				Description: "This is a sample quiz",
				Start_date:  time.Now(),
				End_date:    time.Now().AddDate(0, 0, 7),
			},
			{
				User_id:     1,
				Title:       "Sample Quiz 2",
				Description: "Another sample quiz",
				Start_date:  time.Now(),
				End_date:    time.Now().AddDate(0, 0, 7),
			},
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/quizzes?search=query&page=1&pageSize=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("search", "page", "pageSize")
		c.SetParamValues("query", "1", "10")

		count := int64(10)
		mockRepo.On("GetAllQuiz", 1, 10, "query").Return(fakeQuizData, count, 0)

		err := quizController.GetAllQuiz()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"Succes Get All Quiz\"")

		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("InvalidPageGetAllQuiz", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/quizzes?search=query&page=abc&pageSize=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("search", "page", "pageSize")
		c.SetParamValues("query", "abc", "10")

		fakeQuizData := []model.Quiz{
			{
				User_id:     1,
				Title:       "Sample Quiz 1",
				Description: "This is a sample quiz",
				Start_date:  time.Now(),
				End_date:    time.Now().AddDate(0, 0, 7),
			},
			{
				User_id:     1,
				Title:       "Sample Quiz 2",
				Description: "Another sample quiz",
				Start_date:  time.Now(),
				End_date:    time.Now().AddDate(0, 0, 7),
			},
		}

		count := int64(10)
		mockRepo.On("GetAllQuiz", 1, 10, "query").Return(fakeQuizData, count, 1)

		err := quizController.GetAllQuiz()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("FailedGetAllQuiz", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/quizzes?search=query&page=1&pageSize=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("search", "page", "pageSize")
		c.SetParamValues("query", "1", "10")

		count := int64(0)
		mockRepo.On("GetAllQuiz", 1, 10, "query").Return(nil, count, 2)

		err := quizController.GetAllQuiz()(c)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		responseBody := rec.Body.String()
		assert.Contains(t, responseBody, "\"message\":\"Failed to Get All Quiz\"")

		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("GetAllQuizNotFound", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/quizzes?search=query&page=1&pageSize=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("search", "page", "pageSize")
		c.SetParamValues("query", "1", "10")

		count := int64(0)
		mockRepo.On("GetAllQuiz", 1, 10, "query").Return(nil, count, 3)

		err := quizController.GetAllQuiz()(c)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		responseBody := rec.Body.String()
		assert.Contains(t, responseBody, "\"message\":\"Quiz not Found\"")

		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

}

func TestGetQuizByID(t *testing.T) {
	t.Run("SuccessfulGetQuizByID", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		fakeQuiz := model.Quiz{
			User_id:     1,
			Title:       "Sample Quiz",
			Description: "This is a sample quiz",
			Start_date:  time.Now(),
			End_date:    time.Now().AddDate(0, 0, 7),
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/quiz/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockRepo.On("GetQuizByID", uint(1)).Return(&fakeQuiz, 0) // Mengembalikan data quiz palsu dengan ID 1

		err := quizController.GetQuizByID()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"Succes Get Quiz\"")

		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("InvalidQuizID", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/quiz/s", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("s")

		err := quizController.GetQuizByID()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"invalid quizId\"")

		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
	

	t.Run("FailedGetQuizByID", func(t *testing.T) {
		mockRepo := new(mocks.QuizInterface)
		mockJWT := mocks.NewJWTInterface(t)

		quizController := controller.NewQuizControllInterface(mockRepo, mockJWT)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/quiz/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockRepo.On("GetQuizByID", uint(1)).Return(nil, 1) 

		err := quizController.GetQuizByID()(c)

		assert.NoError(t, err) 

		assert.Equal(t, http.StatusInternalServerError, rec.Code) 

		responseBody := rec.Body.String()
		assert.Contains(t, responseBody, "\"message\":\"Failed to Get Quiz\"") 

		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}


