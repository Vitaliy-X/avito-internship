package handlers_test

import (
	"avito-internship/internal/models"
	"avito-internship/internal/transport/handlers"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockPVZService struct {
	mock.Mock
}

func (m *mockPVZService) CreatePVZ(pvz models.PVZ) (*models.PVZ, error) {
	args := m.Called(pvz)
	if pvzPtr := args.Get(0); pvzPtr != nil {
		return pvzPtr.(*models.PVZ), args.Error(1)
	}
	return nil, args.Error(1)
}

func createPVZHandler(service *mockPVZService) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "moderator" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Only moderators can create PVZ"})
			return
		}

		var req handlers.CreatePVZRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		pvz := models.PVZ{
			ID:               req.ID,
			City:             req.City,
			RegistrationDate: req.RegistrationDate,
		}

		if service == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Service unavailable"})
			return
		}

		result, err := service.CreatePVZ(pvz)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func setupRouterPVZ(handler gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		role := c.GetHeader("Role")
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})

	r.POST("/pvz", func(c *gin.Context) {
		if handler != nil {
			handler(c)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "handler not set"})
		}
	})

	return r
}

func TestCreatePVZ_Success(t *testing.T) {
	mockSvc := new(mockPVZService)

	pvz := models.PVZ{
		ID:               "11111111-1111-1111-1111-111111111111",
		City:             "Москва",
		RegistrationDate: time.Date(2025, 4, 24, 18, 0, 0, 0, time.UTC),
	}

	mockSvc.On("CreatePVZ", mock.MatchedBy(func(input models.PVZ) bool {
		return input.ID == pvz.ID && input.City == pvz.City && input.RegistrationDate.Equal(pvz.RegistrationDate)
	})).Return(&pvz, nil)

	router := setupRouterPVZ(createPVZHandler(mockSvc))

	body, _ := json.Marshal(gin.H{
		"id":               pvz.ID,
		"registrationDate": pvz.RegistrationDate,
		"city":             pvz.City,
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Role", "moderator")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.PVZ
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Москва", response.City)

	mockSvc.AssertExpectations(t)
}

func TestCreatePVZ_InvalidInput(t *testing.T) {
	mockSvc := new(mockPVZService)
	router := setupRouterPVZ(createPVZHandler(mockSvc))

	body := `{"id":"not-a-uuid","city":"Москва"}`

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Role", "moderator")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestCreatePVZ_Forbidden(t *testing.T) {
	mockSvc := new(mockPVZService)
	router := setupRouterPVZ(createPVZHandler(mockSvc))

	body := `{"id":"11111111-1111-1111-1111-111111111111","registrationDate":"2025-04-24T18:00:00Z","city":"Москва"}`
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(body))
	req.Header.Set("Role", "employee")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Only moderators can create PVZ")
}

func TestCreatePVZ_CityNotAllowed(t *testing.T) {
	mockSvc := new(mockPVZService)

	pvz := models.PVZ{
		ID:               "11111111-1111-1111-1111-111111111111",
		City:             "Новосибирск",
		RegistrationDate: time.Date(2025, 4, 24, 18, 0, 0, 0, time.UTC),
	}

	mockSvc.On("CreatePVZ", mock.Anything).Return(nil, errors.New("city not allowed"))

	router := setupRouterPVZ(createPVZHandler(mockSvc))

	body, _ := json.Marshal(gin.H{
		"id":               pvz.ID,
		"registrationDate": pvz.RegistrationDate,
		"city":             pvz.City,
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
	req.Header.Set("Role", "moderator")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "city not allowed")
}
