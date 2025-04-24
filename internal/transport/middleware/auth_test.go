package middleware_test

import (
	"avito-internship/internal/transport/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		role := c.GetString("role")
		c.JSON(http.StatusOK, gin.H{"role": role})
	})
	return router
}

func TestAuthMiddleware_ValidModerator(t *testing.T) {
	router := setupRouter()

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", "Bearer moderator")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "moderator")
}

func TestAuthMiddleware_ValidEmployee(t *testing.T) {
	router := setupRouter()

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", "Bearer employee")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "employee")
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	router := setupRouter()

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, request)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Missing token")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := setupRouter()

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", "Bearer noname")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusForbidden, response.Code)
	assert.Contains(t, response.Body.String(), "Invalid role")
}
