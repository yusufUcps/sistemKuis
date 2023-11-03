package controller_test

import (
	"net/http"
	"net/http/httptest"
	"quiz/mocks"
	"quiz/model"
	"strings"

	"testing"

	"quiz/controller"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertQuestion(t *testing.T) {
	// Create mock objects
	mockRepo := new(mocks.QuestionsInterface)
	mockJWT := new(mocks.JWTInterface)
	mockOpenAI := new(mocks.OpenAiInterface)

	qc := controller.NewQuestionsControllInterface(mockRepo, mockJWT, mockOpenAI)

	e := echo.New()

	// Mock the request payload
	reqPayload := `{"quiz_id": 1, "question": "What is 1+1?", "options": []}`

	req := httptest.NewRequest(http.MethodPost, "/insert-question", strings.NewReader(reqPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock the behavior of your repository and JWT for successful insertion
	mockRepo.On("InsertQuestion", mock.Anything).Return(&model.Questions{}, 0)

	err := qc.InsertQuestion()(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	responseBody := rec.Body.String()
	t.Log("Response Body:", responseBody)

	assert.Contains(t, responseBody, "\"message\":\"Succes create Question\"")

	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)

}
func TestGetAllQuestionsQuiz(t *testing.T) {
    t.Run("SuccessfulGetAllQuestionsQuiz", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.QuestionsInterface)
		mockJWT := new(mocks.JWTInterface)
		mockOpenAI := new(mocks.OpenAiInterface)

		var coba  []model.Questions

		qc := controller.NewQuestionsControllInterface(mockRepo, mockJWT, mockOpenAI)

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-all-questions-quiz?quizId=1&page=1&pageSize=10", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        mockQuizID := 1 

        c.SetPath("/get-all-questions-quiz")
		
        count := int64(1)
 
        mockRepo.On("GetAllQuestionsFromQuiz", 1, 10, uint(mockQuizID)).Return(coba, count, 0)

        err := qc.GetAllQuestionsQuiz()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })

}

func TestGetQuetionByID(t *testing.T) {
    t.Run("SuccessfulGetQuetionByID", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.QuestionsInterface)
		mockJWT := new(mocks.JWTInterface)
		mockOpenAI := new(mocks.OpenAiInterface)

		var coba  model.Questions

		qc := controller.NewQuestionsControllInterface(mockRepo, mockJWT, mockOpenAI)

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-question/1", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues("1")

        mockRepo.On("GetQuestionByID", uint(1)).Return(&coba, 0)

        // Call the GetQuetionByID endpoint
        err := qc.GetQuetionByID()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })

}

func TestUpdateQuestion(t *testing.T) {
    t.Run("SuccessfulUpdateQuestion", func(t *testing.T) {
        mockRepo := new(mocks.QuestionsInterface)
		mockJWT := new(mocks.JWTInterface)
		mockOpenAI := new(mocks.OpenAiInterface)

		var coba  model.Questions

		qc := controller.NewQuestionsControllInterface(mockRepo, mockJWT, mockOpenAI)

        fakeUserID := uint(1)

        userID := uint(1)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            t.Errorf("token signing error: %s", err)
        }

        reqPayload := `{"quiz_id": 1, "question": "What is 1+1?", "options": []}`

        e := echo.New()
        req := httptest.NewRequest(http.MethodPut, "/update-question/1", strings.NewReader(reqPayload))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)
        req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
        c.SetParamNames("id")
        c.SetParamValues("1")

        mockJWT.On("ExtractToken", mock.Anything).Return(uint(fakeUserID))

        mockRepo.On("UpdateQuestion", mock.Anything, uint(fakeUserID)).Return(&coba, 0)

        err = qc.UpdateQuestion()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
        mockJWT.AssertExpectations(t)
    })

}

func TestDeleteQuestion(t *testing.T) {
    t.Run("SuccessfulDeleteQuestion", func(t *testing.T) {
        mockRepo := new(mocks.QuestionsInterface)
		mockJWT := new(mocks.JWTInterface)
		mockOpenAI := new(mocks.OpenAiInterface)

		qc := controller.NewQuestionsControllInterface(mockRepo, mockJWT, mockOpenAI)

        fakeUserID := uint(1)

        userID := uint(1)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            t.Errorf("token signing error: %s", err)
        }

        e := echo.New()
        req := httptest.NewRequest(http.MethodDelete, "/delete-question/1", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)
        req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
        c.SetParamNames("id")
        c.SetParamValues("1")

        mockJWT.On("ExtractToken", mock.Anything).Return(uint(fakeUserID))

        mockRepo.On("DeleteQuestion", uint(1), uint(fakeUserID)).Return(0)

        err = qc.DeleteQuestion()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
        mockJWT.AssertExpectations(t)
    })

}

func TestGenerateQuestion(t *testing.T) {
	t.Run("SuccessfulGenerateQuestion", func(t *testing.T) {
		mockRepo := new(mocks.QuestionsInterface)
		mockJWT := new(mocks.JWTInterface)
		mockOpenAI := new(mocks.OpenAiInterface)

		qc := controller.NewQuestionsControllInterface(mockRepo, mockJWT, mockOpenAI)

		var ques []model.Questions

        question1 := model.Questions{
            Quiz_id:    2,
            Question:  "What is Golang?",

        }

        question2 := model.Questions{
            Quiz_id:    2,
            Question:  "How does Golang handle concurrency?",

        }

        ques = append(ques, question1, question2)

		mockOpenAI.On("GenerateQuestions", mock.Anything, mock.Anything).Return(ques).Once()

		requestPayload := `{
			"quiz_id": 7,
			"description": "Your description"
		}`

		req := httptest.NewRequest(http.MethodPost, "/your-endpoint", strings.NewReader(requestPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

        mockRepo.On("InsertGenerateQuestion", ques).Return(ques,0)

		err := qc.GenerateQuestion()(c)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"Succes Generate Question\"")

        mockRepo.AssertExpectations(t)
        mockJWT.AssertExpectations(t)
		mockOpenAI.AssertExpectations(t)
	})

    t.Run("dailedGenerateQuestion", func(t *testing.T) {
		mockRepo := new(mocks.QuestionsInterface)
		mockJWT := new(mocks.JWTInterface)
		mockOpenAI := new(mocks.OpenAiInterface)

		qc := controller.NewQuestionsControllInterface(mockRepo, mockJWT, mockOpenAI)

		var ques []model.Questions


		mockOpenAI.On("GenerateQuestions", mock.Anything, mock.Anything).Return(ques).Once()

		requestPayload := `{
			"quiz_id": 7,
			"description": "Your description"
		}`

		req := httptest.NewRequest(http.MethodPost, "/your-endpoint", strings.NewReader(requestPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

        mockRepo.On("InsertGenerateQuestion", ques).Return(ques,0)

		err := qc.GenerateQuestion()(c)

		assert.Nil(t, err)
        assert.Equal(t, http.StatusInternalServerError, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"Failed to Generate Qusetions\"")

        mockJWT.AssertExpectations(t)
		mockOpenAI.AssertExpectations(t)
	})
}

