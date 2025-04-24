package services_test

import (
	"avito-internship/internal/models"
	"avito-internship/internal/services"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePVZ_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvz := models.PVZ{
		ID:               "test-id-1",
		RegistrationDate: time.Now(),
		City:             "Москва",
	}

	mock.ExpectQuery(`INSERT INTO pvz`).
		WithArgs(pvz.ID, pvz.RegistrationDate, pvz.City).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow(pvz.ID, pvz.RegistrationDate, pvz.City))

	answer, err := services.CreatePVZ(db, pvz)
	assert.NoError(t, err)
	assert.Equal(t, pvz.ID, answer.ID)
	assert.Equal(t, pvz.City, answer.City)
}

func TestCreatePVZ_CityNotAllowed(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvz := models.PVZ{
		ID:               "test-id-1",
		RegistrationDate: time.Now(),
		City:             "Ростов",
	}

	answer, err := services.CreatePVZ(db, pvz)
	assert.Nil(t, answer)
	assert.Error(t, err)
	assert.Equal(t, "city not allowed", err.Error())
}

func TestCreatePVZ_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pvz := models.PVZ{
		ID:               "test-id-1",
		RegistrationDate: time.Now(),
		City:             "Москва",
	}

	mock.ExpectQuery(`INSERT INTO pvz`).
		WithArgs(pvz.ID, pvz.RegistrationDate, pvz.City).
		WillReturnError(errors.New("insert failed"))

	answer, err := services.CreatePVZ(db, pvz)
	assert.Nil(t, answer)
	assert.Error(t, err)
	assert.Equal(t, "insert failed", err.Error())
}

func TestGetPVZList_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	page := 1
	limit := 5
	offset := 0

	mock.ExpectQuery(`SELECT id, registration_date, city FROM pvz`).
		WithArgs(startDate, endDate, limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow("pvz-id", time.Now(), "Москва"))

	mock.ExpectQuery(`SELECT id, date_time, status FROM receptions`).
		WithArgs("pvz-id", startDate, endDate).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "status"}).
			AddRow("rec-id", time.Now(), "in_progress"))

	mock.ExpectQuery(`SELECT id, date_time, type FROM products`).
		WithArgs("rec-id").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "type"}).
			AddRow("prod-id", time.Now(), "электроника"))

	answer, err := services.GetPVZList(db, &startDate, &endDate, page, limit)
	assert.NoError(t, err)
	assert.Len(t, answer, 1)
	assert.Equal(t, "Москва", answer[0].PVZ.City)
	assert.Len(t, answer[0].Receptions, 1)
	assert.Len(t, answer[0].Receptions[0].Products, 1)
}

func TestGetPVZList_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	page := 1
	limit := 5
	offset := 0

	mock.ExpectQuery(`SELECT id, registration_date, city FROM pvz`).
		WithArgs(startDate, endDate, limit, offset).
		WillReturnError(errors.New("query error"))

	answer, err := services.GetPVZList(db, &startDate, &endDate, page, limit)
	assert.Nil(t, answer)
	assert.Error(t, err)
	assert.Equal(t, "query error", err.Error())
}
