package services_test

import (
	"avito-internship/internal/services"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateReception_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvzID := "pvz-123"

	mock.ExpectQuery(`SELECT id FROM receptions`).
		WithArgs(pvzID).
		WillReturnError(sql.ErrNoRows)

	now := time.Now()
	mock.ExpectQuery(`INSERT INTO receptions`).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "status"}).
			AddRow("rec-id", now, "in_progress"))

	r, err := services.CreateReception(db, pvzID)
	assert.NoError(t, err)
	assert.Equal(t, "in_progress", r.Status)
	assert.Equal(t, pvzID, r.PVZID)
	assert.Equal(t, "rec-id", r.ID)
}

func TestCreateReception_AlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvzID := "pvz-123"

	mock.ExpectQuery(`SELECT id FROM receptions`).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("rec-1"))

	r, err := services.CreateReception(db, pvzID)
	assert.Nil(t, r)
	assert.EqualError(t, err, "already an open reception")
}

func TestCreateReception_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvzID := "pvz-123"

	mock.ExpectQuery(`SELECT id FROM receptions`).
		WithArgs(pvzID).
		WillReturnError(errors.New("db failure"))

	r, err := services.CreateReception(db, pvzID)
	assert.Nil(t, r)
	assert.EqualError(t, err, "db failure")
}

func TestCloseLastReception_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvzID := "pvz-123"

	mock.ExpectExec(`UPDATE receptions`).
		WithArgs(pvzID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	err = services.CloseLastReception(db, pvzID)
	assert.NoError(t, err)
}

func TestCloseLastReception_NoActive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvzID := "pvz-123"

	mock.ExpectExec(`UPDATE receptions`).
		WithArgs(pvzID).
		WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected

	err = services.CloseLastReception(db, pvzID)
	assert.EqualError(t, err, "no active reception")
}

func TestCloseLastReception_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvzID := "pvz-123"

	mock.ExpectExec(`UPDATE receptions`).
		WithArgs(pvzID).
		WillReturnError(errors.New("update failed"))

	err = services.CloseLastReception(db, pvzID)
	assert.EqualError(t, err, "update failed")
}
