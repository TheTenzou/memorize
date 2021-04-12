package controllers

import (
	"bytes"
	"encoding/json"
	"memorize/mocks"
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup(test *testing.T) {
	gin.SetMode(gin.TestMode)

	test.Run("Login and Password Required", func(test *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*models.User"),
		).Return(nil)

		recorder := httptest.NewRecorder()

		router := gin.Default()

		NewController(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		requestBody, err := json.Marshal(gin.H{
			"login": "",
		})

		assert.NoError(test, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(test, 400, recorder.Code)
		mockUserService.AssertNotCalled(test, "Signup")
	})

	test.Run("Password to short", func(test *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*models.User"),
		).Return(nil)

		recorder := httptest.NewRecorder()

		router := gin.Default()

		NewController(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		requestBody, err := json.Marshal(gin.H{
			"login":    "alice",
			"password": "pass",
		})

		assert.NoError(test, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(test, 400, recorder.Code)
		mockUserService.AssertNotCalled(test, "Signup")
	})

	test.Run("Password to long", func(test *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*models.User"),
		).Return(nil)

		recorder := httptest.NewRecorder()

		router := gin.Default()

		NewController(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		requestBody, err := json.Marshal(gin.H{
			"login":    "alice",
			"password": "passsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
		})

		assert.NoError(test, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(test, 400, recorder.Code)
		mockUserService.AssertNotCalled(test, "Signup")
	})

	test.Run("Error returnd from userService", func(test *testing.T) {
		user := &models.User{
			Login:    "alice",
			Password: "alicePassword",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*models.User"),
		).Return(apperrors.NewConflict("User Already Exits", user.Login))

		recorder := httptest.NewRecorder()

		router := gin.Default()

		NewController(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		requestBody, err := json.Marshal(gin.H{
			"login":    user.Login,
			"password": user.Password,
		})

		assert.NoError(test, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(test, 409, recorder.Code)
		mockUserService.AssertExpectations(test)
	})

	test.Run("Successful Token Creation", func(test *testing.T) {
		user := &models.User{
			Login:    "alice",
			Password: "alicePassword",
		}

		mockTokenResponce := &models.TokenPair{
			IDToken:      "idToken",
			RefreshToken: "refreshTOken",
		}

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*context.emptyCtx"),
			user,
		).Return(nil)
		mockTokenService.On(
			"NewPairFromUser",
			mock.AnythingOfType("*context.emptyCtx"),
			user,
			"",
		).Return(mockTokenResponce, nil)

		recorder := httptest.NewRecorder()

		router := gin.Default()

		NewController(&Config{
			Router:       router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		requestBody, err := json.Marshal(gin.H{
			"login":    user.Login,
			"password": user.Password,
		})
		assert.NoError(test, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenResponce,
		})
		assert.NoError(test, err)

		assert.Equal(test, http.StatusCreated, recorder.Code)
		assert.Equal(test, respBody, recorder.Body.Bytes())

		mockUserService.AssertExpectations(test)
		mockTokenService.AssertExpectations(test)
	})

	test.Run("Failed Token Creation", func(test *testing.T) {
		user := &models.User{
			Login:    "alice",
			Password: "avalidPassword",
		}

		mockErrorResponse := apperrors.NewInternal()

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*context.emptyCtx"),
			user,
		).Return(nil)
		mockTokenService.On(
			"NewPairFromUser",
			mock.AnythingOfType("*context.emptyCtx"),
			user,
			"",
		).Return(nil, mockErrorResponse)

		recorder := httptest.NewRecorder()

		router := gin.Default()

		NewController(&Config{
			Router:       router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		requestBody, err := json.Marshal(gin.H{
			"login":    user.Login,
			"password": user.Password,
		})
		assert.NoError(test, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockErrorResponse,
		})
		assert.NoError(test, err)

		assert.Equal(test, mockErrorResponse.Status(), recorder.Code)
		assert.Equal(test, respBody, recorder.Body.Bytes())

		mockUserService.AssertExpectations(test)
		mockTokenService.AssertExpectations(test)
	})
}
