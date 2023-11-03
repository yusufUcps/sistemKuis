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


func TestInsertAnswer(t *testing.T) {
    mockRepo := new(mocks.HistoryInterface)
    mockJWT := new(mocks.JWTInterface)
    mockExport := new(mocks.ExportInterface)

    qc := controller.NewHistoryControllInterface(mockRepo, mockJWT, mockExport)

    userID := uint(1)
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["id"] = userID
    tokenString, err := token.SignedString([]byte("s3cr3t"))
    if err != nil {
        t.Errorf("token signing error: %s", err)
    }

    e := echo.New()

    reqPayload := `[{"question_id": 1, "option_id": 2}]`

    req := httptest.NewRequest(http.MethodPost, "/?quizId=1", strings.NewReader(reqPayload))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)

    c.Set("user", &jwt.Token{Valid: true, Claims: token.Claims})

    mockRepo.On("AnswersInsert", mock.Anything, userID, uint(1)).Return(&model.HistoryScore{}, 0)
    mockJWT.On("ExtractToken", mock.Anything).Return(userID)

    err = qc.Answering()(c)

    assert.Nil(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)

    responseBody := rec.Body.String()
    t.Log("Response Body:", responseBody)

    assert.Contains(t, responseBody, "\"message\":\"Succes save asswer\"")

    mockRepo.AssertExpectations(t)
    mockJWT.AssertExpectations(t)
}

func TestGetAllMyHistoryScore(t *testing.T) {
    t.Run("SuccessfulGetAllMyHistoryScore", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.HistoryInterface)
    	mockJWT := new(mocks.JWTInterface)
    	mockExport := new(mocks.ExportInterface)

    	qc := controller.NewHistoryControllInterface(mockRepo, mockJWT, mockExport)
		userID := uint(1)
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID
		tokenString, err := token.SignedString([]byte("s3cr3t"))
    if err != nil {
        t.Errorf("token signing error: %s", err)
    }

		var coba []model.HistoryScore

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-all-history", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		c.Set("user", &jwt.Token{Valid: true, Claims: token.Claims})
		mockJWT.On("ExtractToken", mock.Anything).Return(userID)

        c.SetPath("/get-all-questions-quiz")

        mockRepo.On("GetAllMyHistoryScore", 1,10,uint(1),"").Return(coba,int64(1),0)

        err = qc.GetAllMyHistoryScore()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })

}

func TestGetAllMyHistoryMyquiz(t *testing.T) {
    t.Run("SuccessfulGetAllMyHistoryMyquiz", func(t *testing.T) {
        
        mockRepo := new(mocks.HistoryInterface)
    	mockJWT := new(mocks.JWTInterface)
    	mockExport := new(mocks.ExportInterface)

    	qc := controller.NewHistoryControllInterface(mockRepo, mockJWT, mockExport)
		userID := uint(1)
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID
		tokenString, err := token.SignedString([]byte("s3cr3t"))
    	if err != nil {
        	t.Errorf("token signing error: %s", err)
    	}

		var coba []model.HistoryScore

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-all-history?quizId=1", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		c.Set("user", &jwt.Token{Valid: true, Claims: token.Claims})
		mockJWT.On("ExtractToken", mock.Anything).Return(userID)

        c.SetPath("/get-all-questions-quiz")

        mockRepo.On("GetAllHistoryScoreMyQuiz", 1,10,uint(1),"",uint(1)).Return(coba,int64(1),0)

        err = qc.GetAllHistoryScoreMyQuiz()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })

}

func TestGetHistoryScoreById(t *testing.T) {
    t.Run("SuccessfulTestGetHistoryScoreById", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.HistoryInterface)
    	mockJWT := new(mocks.JWTInterface)
    	mockExport := new(mocks.ExportInterface)

    	qc := controller.NewHistoryControllInterface(mockRepo, mockJWT, mockExport)
		userID := uint(1)
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID
		tokenString, err := token.SignedString([]byte("s3cr3t"))
    	if err != nil {
        	t.Errorf("token signing error: %s", err)
    	}

		var coba model.HistoryScore

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-question/1", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		c.Set("user", &jwt.Token{Valid: true, Claims: token.Claims})
		mockJWT.On("ExtractToken", mock.Anything).Return(userID)
        c.SetParamNames("id")
        c.SetParamValues("1")

        mockRepo.On("GetHistoryScoreById", uint(1), uint(1)).Return(&coba, 0)

        // Call the GetQuetionByID endpoint
        err = qc.GetHistoryScoreById()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })

}

func TestGetAllHistoryAnswer(t *testing.T) {
    t.Run("SuccessfulGetAllHistoryAnswer", func(t *testing.T) {
        
        mockRepo := new(mocks.HistoryInterface)
    	mockJWT := new(mocks.JWTInterface)
    	mockExport := new(mocks.ExportInterface)

    	qc := controller.NewHistoryControllInterface(mockRepo, mockJWT, mockExport)
		userID := uint(1)
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID
		tokenString, err := token.SignedString([]byte("s3cr3t"))
    	if err != nil {
        	t.Errorf("token signing error: %s", err)
    	}

		var coba []model.HistoryAnswers

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-all-history?historyId=1", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		c.Set("user", &jwt.Token{Valid: true, Claims: token.Claims})
		mockJWT.On("ExtractToken", mock.Anything).Return(userID)

        c.SetPath("/get-all-questions-quiz")

        mockRepo.On("GetAllHistoryAnswer", 1,10,uint(1),uint(1)).Return(coba,int64(1),0)

        err = qc.GetAllHistoryAnswer()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })

}

func TestExAllMyHistoryScore(t *testing.T) {
    t.Run("SuccessfulExportMyHistoryScore", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.HistoryInterface)
    	mockJWT := new(mocks.JWTInterface)
    	mockExport := new(mocks.ExportInterface)

    	qc := controller.NewHistoryControllInterface(mockRepo, mockJWT, mockExport)
		userID := uint(1)
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID
		tokenString, err := token.SignedString([]byte("s3cr3t"))
    if err != nil {
        t.Errorf("token signing error: %s", err)
    }
		var coba []model.HistoryScore

		var ex []model.MyHistoryScoreRes

		var res model.ExportRes

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-all-history", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		c.Set("user", &jwt.Token{Valid: true, Claims: token.Claims})
		mockJWT.On("ExtractToken", mock.Anything).Return(userID)

        c.SetPath("/get-all-questions-quiz")

        mockRepo.On("ExMyHistoryScore", uint(1)).Return(coba,0)
		mockExport.On("ExportMyHistoryScore", ex, uint(1)).Return(&res,nil)

        err = qc.ExportMyHistoryScore()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })
}

func TestExportHistoryScoreMyQuiz(t *testing.T) {
    t.Run("SuccessfulExportHistoryScoreMyQuiz", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.HistoryInterface)
    	mockJWT := new(mocks.JWTInterface)
    	mockExport := new(mocks.ExportInterface)

    	qc := controller.NewHistoryControllInterface(mockRepo, mockJWT, mockExport)
		userID := uint(1)
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID
		tokenString, err := token.SignedString([]byte("s3cr3t"))
    if err != nil {
        t.Errorf("token signing error: %s", err)
    }
		var coba []model.HistoryScore

		var ex []model.HistoryScoreMyQuizRes

		var res model.ExportRes

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-all-history?quizId=1", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		c.Set("user", &jwt.Token{Valid: true, Claims: token.Claims})
		mockJWT.On("ExtractToken", mock.Anything).Return(userID)

        c.SetPath("/get-all-questions-quiz")

        mockRepo.On("ExHistoryScoreMyQuiz", uint(1),uint(1)).Return(coba,0)
		mockExport.On("ExportHistoryScoreMyQuiz", ex, uint(1)).Return(&res,nil)

        err = qc.ExportHistoryScoreMyQuiz()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })
}

func TestExportHistoryAnswer(t *testing.T) {
    t.Run("SuccessfulExportHistoryAnswer", func(t *testing.T) {
        // Create mock objects and fake variables
        mockRepo := new(mocks.HistoryInterface)
    	mockJWT := new(mocks.JWTInterface)
    	mockExport := new(mocks.ExportInterface)

    	qc := controller.NewHistoryControllInterface(mockRepo, mockJWT, mockExport)
		userID := uint(1)
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID
		tokenString, err := token.SignedString([]byte("s3cr3t"))
    if err != nil {
        t.Errorf("token signing error: %s", err)
    }
		var coba []model.HistoryAnswers

		var ex []model.HistoryAnswersRes

		var res model.ExportRes

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/get-all-history?historyId=1", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		c.Set("user", &jwt.Token{Valid: true, Claims: token.Claims})
		mockJWT.On("ExtractToken", mock.Anything).Return(userID)

        c.SetPath("/get-all-questions-quiz")

        mockRepo.On("ExHistoryAnswer", uint(1),uint(1)).Return(coba,0)
		mockExport.On("ExportHistoryAnswer", ex, uint(1)).Return(&res,nil)

        err = qc.ExportHistoryAnswer()(c)

        assert.Nil(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockRepo.AssertExpectations(t)
    })
}