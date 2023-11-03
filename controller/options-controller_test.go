package controller_test

import (
	"net/http"
	"net/http/httptest"
	"quiz/mocks"
	"quiz/model"
	"strings"

	"testing"

	"quiz/controller"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertOptions(t *testing.T) {
	// Create mock objects
	mockRepo := new(mocks.OptionsInterface)
	mockJWT := new(mocks.JWTInterface)


	qc := controller.NewOptionsControllInterface(mockRepo, mockJWT)

	e := echo.New()

	reqPayload := `{"question_id": 1, "value": "What is 1+1?", "ist_right": true}`

	req := httptest.NewRequest(http.MethodPost, "/insert", strings.NewReader(reqPayload))
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