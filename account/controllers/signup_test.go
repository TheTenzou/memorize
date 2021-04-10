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

func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Login and Password Required", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*gin.Context"),
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

		assert.NoError(t, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/api/account/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(t, 400, recorder.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password to short", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*gin.Context"),
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

		assert.NoError(t, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/api/account/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(t, 400, recorder.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password to long", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*gin.Context"),
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

		assert.NoError(t, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/api/account/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(t, 400, recorder.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Error returnd from userService", func(t *testing.T) {
		user := &models.User{
			Login:    "alice",
			Password: "alicePassword",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On(
			"Signup",
			mock.AnythingOfType("*gin.Context"),
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

		assert.NoError(t, err)

		request, err := http.NewRequest(
			http.MethodPost,
			"/api/account/signup",
			bytes.NewBuffer(requestBody),
		)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(t, 409, recorder.Code)
		mockUserService.AssertExpectations(t)
	})
}
