package handlers_test

import (
	"avito-internship/internal/transport/handlers"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/dummyLogin", handlers.DummyLogin)
	return router
}

func TestDummyLogin_ValidModerator(t *testing.T) {
	router := setupRouter()

	body := map[string]string{"role": "moderator"}
	jsonBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"token":"moderator"`)
}

func TestDummyLogin_ValidEmployee(t *testing.T) {
	router := setupRouter()

	body := map[string]string{"role": "employee"}
	jsonBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"token":"employee"`)
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	router := setupRouter()

	body := map[string]string{"role": "admin"}
	jsonBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "Invalid role")
}

func TestDummyLogin_MissingRole(t *testing.T) {
	router := setupRouter()

	body := map[string]string{}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "Invalid role")
}
