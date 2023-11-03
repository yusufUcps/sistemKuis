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
	"github.com/golang-jwt/jwt/v5"
)


func TestRegister(t *testing.T) {
    t.Run("Success", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        user := model.Users{
            Name:  "aku",
            Email: "test@example.com",
        }
        mockRepo.On("Register", mock.Anything).Return(&user, 0)

        e := echo.New()
        req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"aku", "email": "test@example.com", "password": "test123"}`))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        err := userController.Register()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"Succes create account\"")
    })

    t.Run("EmailAlreadyRegistered", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        user := model.Users{
            Name:     "existing",
            Email:    "existing@example.com",
            Password: "existing123",
        }

        mockRepo.On("Register", user).Return(nil, 1)

        e := echo.New()
        req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"existing", "email": "existing@example.com", "password": "existing123"}`))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        err := userController.Register()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusBadRequest, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"Email already registered\"")
    })

    t.Run("InternalServerError", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        user := model.Users{
            Name:     "error",
            Email:    "error@example.com",
            Password: "error123",
        }

        mockRepo.On("Register", user).Return(nil, 2)

        e := echo.New()
        req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"error", "email": "error@example.com", "password": "error123"}`))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        err := userController.Register()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusInternalServerError, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"Failed Register User\"")
    })

    t.Run("InvalidUserInput", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        malformedJSON := `{"name":"John", "email": "test@example.com"`

        e := echo.New()
        req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(malformedJSON))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        err := userController.Register()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusBadRequest, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"invalid user input\"")
    })

	t.Run("InvalidPassword", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        invalidPasswordUser := model.Users{
            Name:     "invalidPassword",
            Email:    "invalid@example.com",
            Password: "short", // Password is too short (invalid)
        }

        mockRepo.On("Register", invalidPasswordUser).Return(nil, 3)

        e := echo.New()
        req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"invalidPassword", "email": "invalid@example.com", "password": "short"}`))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        err := userController.Register()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusBadRequest, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"invalid input password\"")
    })
}

func TestLogin(t *testing.T) {
    t.Run("ValidLogin", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
		mockJWT := mocks.NewJWTInterface(t)

		userController := controller.NewUserControllInterface(mockRepo, mockJWT)

		user := model.Users{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		loginInput := model.Login{
			Email:    "john@example.com",
			Password: "password123",
		}

		mockRepo.On("Login", loginInput.Email, loginInput.Password).Return(&user, 0)

		id := uint(0)
		token := "example_token"
		mockJWT.On("GenerateJWT", id).Return(token).Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email": "john@example.com", "password": "password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := userController.Login()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"Login Success\"")
    })

    t.Run("InvalidUserInputLogin", func(t *testing.T) {
		mockRepo := new(mocks.UsersInterface)
		mockJWT := mocks.NewJWTInterface(t)
	
		userController := controller.NewUserControllInterface(mockRepo, mockJWT)
	
		invalidInput := `{"email": "test@example.com"` 
	
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(invalidInput))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
	
		err := userController.Login()(c)
	
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	
		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)
	
		assert.Contains(t, responseBody, "\"message\":\"invalid user input\"")
    })

    t.Run("InternalServerErrorLogin", func(t *testing.T) {
		mockRepo := new(mocks.UsersInterface)
		mockJWT := mocks.NewJWTInterface(t)
	
		userController := controller.NewUserControllInterface(mockRepo, mockJWT)
	
		loginInput := model.Login{
			Email:    "error@example.com",
			Password: "errorpass",
		}
		mockRepo.On("Login", loginInput.Email, loginInput.Password).Return(nil, 1)
	
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email": "error@example.com", "password": "errorpass"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
	
		err := userController.Login()(c)
	
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	
		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)
	
		assert.Contains(t, responseBody, "\"message\":\"Failed to Login\"")
    })

    t.Run("InvalidCredentialsLogin", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
		mockJWT := mocks.NewJWTInterface(t)

		userController := controller.NewUserControllInterface(mockRepo, mockJWT)

		loginInput := model.Login{
			Email:    "invalid@example.com",
			Password: "invalidpass",
		}
		mockRepo.On("Login", loginInput.Email, loginInput.Password).Return(nil, 2)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email": "invalid@example.com", "password": "invalidpass"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := userController.Login()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"wrong email or password\"")
    })

    t.Run("JWTGenerationFailedLogin", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
		mockJWT := mocks.NewJWTInterface(t)

		userController := controller.NewUserControllInterface(mockRepo, mockJWT)
		
		user := model.Users{
			Name:     "TokenFail",
			Email:    "tokenfail@example.com",
			Password: "tokenfailpass",
		}
		loginInput := model.Login{
			Email:    "tokenfail@example.com",
			Password: "tokenfailpass",
		}
		mockRepo.On("Login", loginInput.Email, loginInput.Password).Return(&user, 0)

		id := uint(0)
		emptyToken := "" 
		mockJWT.On("GenerateJWT", id).Return(emptyToken)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email": "tokenfail@example.com", "password": "tokenfailpass"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := userController.Login()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)

		assert.Contains(t, responseBody, "\"message\":\"cannot process jwt token\"")
    })
}


func TestMyProfile(t *testing.T) {
    t.Run("SuccessfulMyProfile", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        user := model.Users{
            Name:     "aku",
            Email:    "test@example.com",
            Password: "test123",
        }

        userID := uint(0)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            return
        }

        mockJWT.On("ExtractToken", mock.AnythingOfType("*jwt.Token")).Return(userID)

        mockRepo.On("MyProfile", userID).Return(&user, 0)

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)

        req.Header.Set(echo.HeaderAuthorization, "Bearer " + tokenString)

        err = userController.MyProfile()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"success get my profil\"")

        mockJWT.AssertExpectations(t)
        mockRepo.AssertExpectations(t)
    })

    t.Run("InvalidMyProfile", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        userID := uint(0)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            return
        }

        mockJWT.On("ExtractToken", mock.AnythingOfType("*jwt.Token")).Return(userID)

        mockRepo.On("MyProfile", userID).Return(nil, 1)

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)

        req.Header.Set(echo.HeaderAuthorization, "Bearer " + tokenString)

        err = userController.MyProfile()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusInternalServerError, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"Failed to Get MyProfile\"")

        mockJWT.AssertExpectations(t)
        mockRepo.AssertExpectations(t)
    })
}

func TestUpdateMyProfile(t *testing.T) {
    t.Run("SuccessfulUpdateMyProfile", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        uc := controller.NewUserControllInterface(mockRepo, mockJWT)

        user := model.Users{
            Name:     "aku",
            Email:    "test@example.com",
            Password: mock.Anything,
        }

        userID := uint(0)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            return
        }

        mockJWT.On("ExtractToken", mock.AnythingOfType("*jwt.Token")).Return(userID)

        mockRepo.On("UpdateMyProfile", mock.Anything).Return(&user, 0)

        e := echo.New()
        reqPayload := strings.NewReader(`{"Name":"updated name", "Email":"updated@example.com", "Password":"updated123"}`)
        req := httptest.NewRequest(http.MethodPut, "/user", reqPayload)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        req.Header.Set(echo.HeaderAuthorization, "Bearer " + tokenString)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)

        err = uc.UpdateMyProfile()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"success update user profile\"")

        mockJWT.AssertExpectations(t)
        mockRepo.AssertExpectations(t)
    })

    t.Run("UpdateMyProfile_BadRequest", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        uc := controller.NewUserControllInterface(mockRepo, mockJWT)

        userID := uint(1)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            return
        }

        mockJWT.On("ExtractToken", mock.AnythingOfType("*jwt.Token")).Return(userID)

        mockRepo.On("UpdateMyProfile", mock.Anything).Return(nil, 2)

        e := echo.New()
        reqPayload := strings.NewReader(`{"Name":"updated name", "Email":"updated@example.com", "Password":"updated123"}`)
        req := httptest.NewRequest(http.MethodPut, "/user", reqPayload)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        req.Header.Set(echo.HeaderAuthorization, "Bearer " + tokenString)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)

        err = uc.UpdateMyProfile()(c)

        assert.Equal(t, http.StatusBadRequest, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)
        assert.Contains(t, responseBody, "\"message\":\"Email already registered\"")

        mockJWT.AssertExpectations(t)
        mockRepo.AssertExpectations(t)
    })

    t.Run("UpdateMyProfileInternalServerError", func(t *testing.T) {
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        uc := controller.NewUserControllInterface(mockRepo, mockJWT)

        userID := uint(1)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            return
        }

        mockJWT.On("ExtractToken", mock.AnythingOfType("*jwt.Token")).Return(userID)

        mockRepo.On("UpdateMyProfile", mock.Anything).Return(nil, 1)

        e := echo.New()
        reqPayload := strings.NewReader(`{"Name":"updated name", "Email":"updated@example.com", "Password":"updated123"}`)
        req := httptest.NewRequest(http.MethodPut, "/user", reqPayload)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        req.Header.Set(echo.HeaderAuthorization, "Bearer " + tokenString)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)

        err = uc.UpdateMyProfile()(c)

        assert.Equal(t, http.StatusInternalServerError, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        mockJWT.AssertExpectations(t)
        mockRepo.AssertExpectations(t)
    })

	t.Run("UpdateMyProfileInvalidInput", func(t *testing.T) {
		mockRepo := new(mocks.UsersInterface)
		mockJWT := mocks.NewJWTInterface(t)
	
		uc := controller.NewUserControllInterface(mockRepo, mockJWT)
	
		userID := uint(0) 
	
		mockJWT.On("ExtractToken", mock.AnythingOfType("*jwt.Token")).Return(userID)
		
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userID
	
		tokenString, err := token.SignedString([]byte("s3cr3t")) 
		if err != nil {
			t.Errorf("Error signing token: %s", err)
		}
	
		e := echo.New()
		reqPayload := strings.NewReader(`{"Invalid JSON"}`)
		req := httptest.NewRequest(http.MethodPut, "/user", reqPayload)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
	
		err = uc.UpdateMyProfile()(c)
	
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	
		responseBody := rec.Body.String()
		t.Log("Response Body:", responseBody)
	
		assert.Contains(t, responseBody, "invalid user input") 
	})
}

func TestDeleteUser(t *testing.T) {
    t.Run("SuccessfulDeleteUser", func(t *testing.T) {
        
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)
        
        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        userID := uint(1) 

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            t.Errorf("token signing error: %s", err)
        }

        
        mockJWT.On("ExtractToken", mock.Anything).Return(userID)
        mockRepo.On("DeleteUser", userID).Return(0) 

        e := echo.New()
        req := httptest.NewRequest(http.MethodDelete, "/user", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)

        req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)

        
        err = userController.DeleteUser()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"success Delete User\"")

        
        mockJWT.AssertExpectations(t)
        mockRepo.AssertExpectations(t)
    })

    t.Run("FailedDeleteUser", func(t *testing.T) {
        
        mockRepo := new(mocks.UsersInterface)
        mockJWT := mocks.NewJWTInterface(t)

        
        userController := controller.NewUserControllInterface(mockRepo, mockJWT)

        userID := uint(1) 

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = userID

        tokenString, err := token.SignedString([]byte("s3cr3t"))
        if err != nil {
            t.Errorf("token signing error: %s", err)
        }

        
        mockJWT.On("ExtractToken", mock.Anything).Return(userID)
        mockRepo.On("DeleteUser", userID).Return(1) 

        e := echo.New()
        req := httptest.NewRequest(http.MethodDelete, "/user", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", token)

        req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)

        err = userController.DeleteUser()(c)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusInternalServerError, rec.Code)

        responseBody := rec.Body.String()
        t.Log("Response Body:", responseBody)

        assert.Contains(t, responseBody, "\"message\":\"failed to Delete User\"")

        mockJWT.AssertExpectations(t)
        mockRepo.AssertExpectations(t)
    })
}
