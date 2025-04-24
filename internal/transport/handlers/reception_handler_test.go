package handlers_test

import (
	"avito-internship/internal/database"
	"avito-internship/internal/transport/handlers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCreateReceptionHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/receptions", func(c *gin.Context) {
		c.Set("role", "employee")
		handlers.CreateReception(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBufferString("{ invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateReceptionHandler_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/receptions", func(c *gin.Context) {
		c.Set("role", "client")
		handlers.CreateReception(c)
	})

	body := `{"pvzId":"any-id"}`
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestCreateReceptionHandler_RepoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	database.DB = db

	pvzID := "31ae2e29-0460-4748-a9f3-2b5747f78960"

	rows := sqlmock.NewRows([]string{"id"}).AddRow("existing_id")
	mock.ExpectQuery(`SELECT id FROM receptions`).
		WithArgs(pvzID).
		WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/receptions", func(c *gin.Context) {
		c.Set("role", "employee")
		handlers.CreateReception(c)
	})

	body := `{"pvzId":"` + pvzID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "already an open reception", resp["message"])

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseReceptionHandler_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	database.DB = db

	pvzID := "82cc7cda-bd24-468f-b7b7-844d66b6693c"

	mock.ExpectExec(`UPDATE receptions`).
		WithArgs(pvzID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/pvz/:pvzId/close_last_reception", func(c *gin.Context) {
		c.Set("role", "employee")
		handlers.CloseReception(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID+"/close_last_reception", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "reception has been closed", resp["message"])

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseReceptionHandler_NoActiveReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	database.DB = db

	pvzID := "123"
	mock.ExpectExec(`UPDATE receptions`).
		WithArgs(pvzID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/pvz/:pvzId/close_last_reception", func(c *gin.Context) {
		c.Set("role", "employee")
		handlers.CloseReception(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID+"/close_last_reception", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "no active reception", resp["message"])

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseReceptionHandler_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/pvz/:pvzId/close_last_reception", func(c *gin.Context) {
		c.Set("role", "client")
		handlers.CloseReception(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/pvz/111/close_last_reception", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}
