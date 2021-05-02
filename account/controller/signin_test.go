package controller

import (
	"bytes"
	"encoding/json"
	"memorize/mocks"
	"memorize/models"
	"memorize/models/apperrors"
	"memorize/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignin(test *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// setup mock services, gin engine/router, handler layer
	mockUserService := new(mocks.MockUserService)
	mockTokenService := new(mocks.MockTokenService)

	router := gin.Default()

	NewController(&Config{
		Router:       router,
		UserService:  mockUserService,
		TokenService: mockTokenService,
	})

	test.Run("Bad request data", func(test *testing.T) {
		recorder := httptest.NewRecorder()

		requestBody, err := json.Marshal(gin.H{
			"login":    "",
			"password": "short",
		})
		assert.NoError(test, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(requestBody))
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		assert.Equal(test, http.StatusBadRequest, recorder.Code)
		mockUserService.AssertNotCalled(test, "Signin")
		mockTokenService.AssertNotCalled(test, "NewTokensFromUser")
	})

	test.Run("Error Returned from UserService.Signin", func(test *testing.T) {
		login := "alice"
		password := "alicepassword"

		mockUserServiceArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&models.User{
				Login:    login,
				Password: password,
			},
		}

		mockError := apperrors.NewAuthorization("invalid email/password combo")

		mockUserService.On("Signin", mockUserServiceArguments...).Return(nil, mockError)

		recorder := httptest.NewRecorder()

		requestBody, err := json.Marshal(gin.H{
			"login":    login,
			"password": password,
		})
		assert.NoError(test, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(requestBody))
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		mockUserService.AssertCalled(test, "Signin", mockUserServiceArguments...)
		mockTokenService.AssertNotCalled(test, "NewTokensFromUser")
		assert.Equal(test, http.StatusUnauthorized, recorder.Code)
	})

	test.Run("Successful Token Creation", func(test *testing.T) {
		uid, _ := uuid.NewRandom()
		login := "alice"
		password := "alicePassword"
		hashedPassword, _ := service.HashPassword(password)

		mockUserServiceArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&models.User{
				Login:    login,
				Password: password,
			},
		}

		mockUser := &models.User{
			UID:      uid,
			Login:    login,
			Password: hashedPassword,
		}

		mockUserService.On("Signin", mockUserServiceArguments...).Return(mockUser, nil)

		mockTokenServiceArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
			"",
		}

		mockTokenPair := &models.TokenPair{
			AccessToken: models.AccessToken{
				Token: "idToken",
			},
			RefreshToken: models.RefreshToken{
				Token: "refreshToken",
			},
		}

		mockTokenService.On("NewPairFromUser", mockTokenServiceArguments...).Return(mockTokenPair, nil)

		recorder := httptest.NewRecorder()

		requestBody, err := json.Marshal(gin.H{
			"login":    login,
			"password": password,
		})
		assert.NoError(test, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(requestBody))
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		responseBody, err := json.Marshal(gin.H{
			"tokens": mockTokenPair,
		})
		assert.NoError(test, err)

		assert.Equal(test, http.StatusOK, recorder.Code)
		assert.Equal(test, responseBody, recorder.Body.Bytes())

		mockUserService.AssertCalled(test, "Signin", mockUserServiceArguments...)
		mockTokenService.AssertCalled(test, "NewPairFromUser", mockTokenServiceArguments...)
	})

	test.Run("Failed Token Creation", func(test *testing.T) {
		uid, _ := uuid.NewRandom()
		login := "cannotproducetoken"
		password := "cannotproducetoken"
		hashedPassword, _ := service.HashPassword(password)

		mockUserServiceArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&models.User{
				Login:    login,
				Password: password,
			},
		}

		mockUser := &models.User{
			UID:      uid,
			Login:    login,
			Password: hashedPassword,
		}

		mockUserService.On("Signin", mockUserServiceArguments...).Return(mockUser, nil)

		mockTokenServiceArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
			"",
		}

		mockError := apperrors.NewInternal()
		mockTokenService.On("NewPairFromUser", mockTokenServiceArguments...).Return(nil, mockError)
		recorder := httptest.NewRecorder()

		requestBody, err := json.Marshal(gin.H{
			"login":    login,
			"password": password,
		})
		assert.NoError(test, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(requestBody))
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(test, err)

		assert.Equal(test, mockError.Status(), recorder.Code)
		assert.Equal(test, respBody, recorder.Body.Bytes())

		mockUserService.AssertCalled(test, "Signin", mockUserServiceArguments...)
		mockTokenService.AssertCalled(test, "NewPairFromUser", mockTokenServiceArguments...)
	})
}
