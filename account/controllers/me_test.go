package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"memorize/mocks"
	"memorize/models"
	"memorize/models/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &models.User{
			UID:   uid,
			Email: "alice@alice.com",
			Name:  "Alice",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.AnythingOfType("*context.emptyCtx"), uid).Return(mockUserResp, nil)

		recorder := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", &models.User{
				UID: uid,
			})
		})

		NewController(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)

		assert.NoError(t, err)

		router.ServeHTTP(recorder, request)

		respBody, err := json.Marshal(gin.H{
			"user": mockUserResp,
		})

		assert.NoError(t, err)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, respBody, recorder.Body.Bytes())
		mockUserService.AssertExpectations(t)
	})

	t.Run("NoContextUser", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.Anything, mock.Anything).Return(nil, nil)

		recorder := httptest.NewRecorder()

		router := gin.Default()
		NewController(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(recorder, request)

		assert.Equal(t, 500, recorder.Code)
		mockUserService.AssertNotCalled(t, "Get", mock.Anything)
	})

	t.Run("NotFound", func(t *testing.T) {
		uid, _ := uuid.NewRandom()
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down call chain"))

		recorder := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", &models.User{
				UID: uid,
			})
		})

		NewController(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(recorder, request)

		respErr := apperrors.NewNotFound("user", uid.String())

		respBody, err := json.Marshal(gin.H{
			"error": respErr,
		})

		assert.NoError(t, err)

		assert.Equal(t, respErr.Status(), recorder.Code)
		assert.Equal(t, respBody, recorder.Body.Bytes())
		mockUserService.AssertExpectations(t)
	})
}
