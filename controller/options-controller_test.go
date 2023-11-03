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

func TestInsertOptions(t *testing.T) {
	mockRepo := new(mocks.OptionsInterface)
	mockJWT := new(mocks.JWTInterface)


	qc := controller.NewOptionsControllInterface(mockRepo, mockJWT)

	e := echo.New()

	reqPayload := `{"question_id": 1, "value": "What is 1+1?", "ist_right": true}`

	req := httptest.NewRequest(http.MethodPost, "/option", strings.NewReader(reqPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo.On("InsertOption", mock.Anything).Return(&model.Options{}, 0)

	err := qc.InsertOption()(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	responseBody := rec.Body.String()
	t.Log("Response Body:", responseBody)

	assert.Contains(t, responseBody, "\"message\":\"Succes create Question\"")

	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)

}

func TestGetAllOptionsQuiz(t *testing.T) {
    t.Run("SuccessfulGetAllOptionsQuiz", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.OptionsInterface)
		mockJWT := new(mocks.JWTInterface)

		qc := controller.NewOptionsControllInterface(mockRepo, mockJWT)

		var coba []model.Options

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/option?questionId=1", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        c.SetPath("/get-all-questions-quiz")

        mockRepo.On("GetAllOptionsFromQuiz", uint(1)).Return(coba,0)

        err := qc.GetAllOptionsQuiz()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })

}

func TestGetOptionByID(t *testing.T) {
    t.Run("SuccessfuloptionByID", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.OptionsInterface)
		mockJWT := new(mocks.JWTInterface)

		qc := controller.NewOptionsControllInterface(mockRepo, mockJWT)

		var coba model.Options

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/option/1", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues("1")

        mockRepo.On("GetOptionByID", uint(1)).Return(&coba, 0)

        // Call the GetQuetionByID endpoint
        err := qc.GetOptionByID()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })

}

func TestUpdateOption(t *testing.T) {
    t.Run("SuccessfulUpdateOption", func(t *testing.T) {
        mockRepo := new(mocks.OptionsInterface)
		mockJWT := new(mocks.JWTInterface)

		qc := controller.NewOptionsControllInterface(mockRepo, mockJWT)

		var coba model.Options

        fakeUserID := uint(1)

        userID := uint(1)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            t.Errorf("token signing error: %s", err)
        }

        reqPayload := `{"question_id": 1, "value": "What is 1+1?", "ist_right": true}`

        e := echo.New()
        req := httptest.NewRequest(http.MethodPut, "/option/1", strings.NewReader(reqPayload))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)
        req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
        c.SetParamNames("id")
        c.SetParamValues("1")

        mockJWT.On("ExtractToken", mock.Anything).Return(uint(fakeUserID))

        mockRepo.On("UpdateOption", mock.Anything, uint(fakeUserID)).Return(&coba, 0)

        err = qc.UpdateOption()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
        mockJWT.AssertExpectations(t)
    })

}

func TestDeleteOption(t *testing.T) {
    t.Run("SuccessfulDeleteOption", func(t *testing.T) {
		mockRepo := new(mocks.OptionsInterface)
		mockJWT := new(mocks.JWTInterface)

		qc := controller.NewOptionsControllInterface(mockRepo, mockJWT)

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
        req := httptest.NewRequest(http.MethodDelete, "/option/1", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)
        req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
        c.SetParamNames("id")
        c.SetParamValues("1")

        mockJWT.On("ExtractToken", mock.Anything).Return(uint(fakeUserID))

        mockRepo.On("DeleteOption", uint(1), uint(fakeUserID)).Return(0)

        err = qc.DeleteOption()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
        mockJWT.AssertExpectations(t)
    })

}