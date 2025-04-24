package handlers_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito-internship/internal/models"
	"avito-internship/internal/transport/handlers"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) AddProduct(pvzID, typ string) (*models.Product, error) {
	args := m.Called(pvzID, typ)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockService) DeleteLastProduct(pvzID string) error {
	args := m.Called(pvzID)
	return args.Error(0)
}

func setupRouterWithService(service *mockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if role := c.GetHeader("Role"); role != "" {
			c.Set("role", role)
		}
		c.Next()
	})

	r.POST("/products", func(c *gin.Context) {
		handlers.AddProduct(c)
	})
	r.DELETE("/products/:pvzId", func(c *gin.Context) {
		handlers.DeleteLastProduct(c)
	})

	return r
}

func TestAddProduct_InvalidInput(t *testing.T) {
	mockSvc := new(mockService)
	router := setupRouterWithService(mockSvc)

	body := []byte(`{"type":"мебель"}`)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Role", "employee")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid input")
}

func TestAddProduct_Forbidden(t *testing.T) {
	mockSvc := new(mockService)
	router := setupRouterWithService(mockSvc)

	body := []byte(`{"type":"электроника","pvzId":"test-pvz"}`)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Role", "moderator")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.Contains(t, w.Body.String(), "Only employees can add products")
}

func TestDeleteLastProduct_Forbidden(t *testing.T) {
	mockSvc := new(mockService)
	router := setupRouterWithService(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/products/test-pvz", nil)
	req.Header.Set("Role", "moderator")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.Contains(t, w.Body.String(), "Only employees can delete products")
}
