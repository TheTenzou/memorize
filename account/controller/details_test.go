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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDetails(test *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	uid, _ := uuid.NewRandom()
	ctxUser := &models.User{
		UID: uid,
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user", ctxUser)
	})

	mockUserService := new(mocks.MockUserService)

	NewController(&Config{
		Router:      router,
		UserService: mockUserService,
	})

	test.Run("Data binding error", func(test *testing.T) {
		recorder := httptest.NewRecorder()

		reqBody, _ := json.Marshal(gin.H{
			"email": "notanemail",
		})
		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		assert.Equal(test, http.StatusBadRequest, recorder.Code)
		mockUserService.AssertNotCalled(test, "UpdateDetails")
	})

	test.Run("Update success", func(test *testing.T) {
		recorder := httptest.NewRecorder()

		newName := "alice"
		newEmail := "alice@mail.com"
		newWebsite := "https://alice.me"

		requestBody, _ := json.Marshal(gin.H{
			"name":    newName,
			"email":   newEmail,
			"website": newWebsite,
		})

		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")

		userToUpdate := &models.User{
			UID:     ctxUser.UID,
			Name:    newName,
			Email:   newEmail,
			Website: newWebsite,
		}

		updateArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userToUpdate,
		}

		dbImageURL := "https://alice.me/static/696292a38f49.jpg"

		mockUserService.
			On("UpdateDetails", updateArguments...).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*models.User)
				userArg.ImageURL = dbImageURL
			}).
			Return(nil)

		router.ServeHTTP(recorder, request)

		userToUpdate.ImageURL = dbImageURL
		responseBody, _ := json.Marshal(gin.H{
			"user": userToUpdate,
		})

		assert.Equal(test, http.StatusOK, recorder.Code)
		assert.Equal(test, responseBody, recorder.Body.Bytes())
		mockUserService.AssertCalled(test, "UpdateDetails", updateArguments...)
	})

	test.Run("Update failure", func(test *testing.T) {
		recorder := httptest.NewRecorder()

		newName := "alice"
		newEmail := "alice@mail.com"
		newWebsite := "https://alice.me"

		requestBody, _ := json.Marshal(gin.H{
			"name":    newName,
			"email":   newEmail,
			"website": newWebsite,
		})

		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")

		userToUpdate := &models.User{
			UID:     ctxUser.UID,
			Name:    newName,
			Email:   newEmail,
			Website: newWebsite,
		}

		updateArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userToUpdate,
		}

		mockError := apperrors.NewInternal()

		mockUserService.
			On("UpdateDetails", updateArguments...).
			Return(mockError)

		router.ServeHTTP(recorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(test, mockError.Status(), recorder.Code)
		assert.Equal(test, responseBody, recorder.Body.Bytes())
		mockUserService.AssertCalled(test, "UpdateDetails", updateArguments...)
	})
}
