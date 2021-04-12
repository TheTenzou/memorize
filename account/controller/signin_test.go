package controller

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

		// create a request body with invalid fields
		reqBody, err := json.Marshal(gin.H{
			"login":    "",
			"password": "short",
		})
		assert.NoError(test, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
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

		mockUserServiceArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&models.User{
				Login:    login,
				Password: password,
			},
		}

		mockError := apperrors.NewAuthorization("invalid email/password combo")

		mockUserService.On("Signin", mockUserServiceArgs...).Return(mockError)

		recorder := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"login":    login,
			"password": password,
		})
		assert.NoError(test, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(test, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		mockUserService.AssertCalled(test, "Signin", mockUserServiceArgs...)
		mockTokenService.AssertNotCalled(test, "NewTokensFromUser")
		assert.Equal(test, http.StatusUnauthorized, recorder.Code)
	})

	test.Run("Successful Token Creation", func(t *testing.T) {
		login := "alice"
		password := "alicePassword"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&models.User{
				Login:    login,
				Password: password,
			},
		}

		mockUserService.On("Signin", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&models.User{
				Login:    login,
				Password: password,
			},
			"",
		}

		mockTokenPair := &models.TokenPair{
			IDToken:      "idToken",
			RefreshToken: "refreshToken",
		}

		mockTokenService.On("NewPairFromUser", mockTSArgs...).Return(mockTokenPair, nil)

		recorder := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"login":    login,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenPair,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, respBody, recorder.Body.Bytes())

		mockUserService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})

	test.Run("Failed Token Creation", func(t *testing.T) {
		login := "cannotproducetoken"
		password := "cannotproducetoken"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&models.User{
				Login:    login,
				Password: password,
			},
		}

		mockUserService.On("Signin", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&models.User{
				Login:    login,
				Password: password,
			},
			"",
		}

		mockError := apperrors.NewInternal()
		mockTokenService.On("NewPairFromUser", mockTSArgs...).Return(nil, mockError)
		recorder := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"login":    login,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), recorder.Code)
		assert.Equal(t, respBody, recorder.Body.Bytes())

		mockUserService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})
}
