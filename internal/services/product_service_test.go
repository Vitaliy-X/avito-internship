package services_test

import (
	"avito-internship/internal/services"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestAddProduct_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := sqlx.NewDb(sqlDB, "postgres")
	pvzID := "pvz-1"
	receptionID := "rec-1"
	now := time.Now()

	mock.ExpectQuery(`(?i)SELECT id FROM receptions.*status = 'in_progress'`).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionID))

	mock.ExpectQuery(`(?i)INSERT INTO products`).
		WithArgs("обувь", receptionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time"}).AddRow("prod-1", now))

	product, err := services.AddProduct(db.DB, pvzID, "обувь")
	require.NoError(t, err)
	require.NotNil(t, product)
	require.Equal(t, "обувь", product.Type)
	require.Equal(t, receptionID, product.ReceptionID)
	require.Equal(t, "prod-1", product.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddProduct_NoReception(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := sqlx.NewDb(sqlDB, "postgres")
	pvzID := "pvz-1"

	mock.ExpectQuery(`(?i)SELECT id FROM receptions.*status = 'in_progress'`).
		WithArgs(pvzID).
		WillReturnError(sql.ErrNoRows)

	product, err := services.AddProduct(db.DB, pvzID, "обувь")
	require.Error(t, err)
	require.Nil(t, product)
	require.EqualError(t, err, "no active product")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLastProduct_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := sqlx.NewDb(sqlDB, "postgres")
	pvzID := "pvz-1"
	productID := "prod-1"

	mock.ExpectQuery(`(?i)SELECT p.id FROM products`).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(productID))

	mock.ExpectExec(`(?i)DELETE FROM products`).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = services.DeleteLastProduct(db.DB, pvzID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLastProduct_NoProduct(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := sqlx.NewDb(sqlDB, "postgres")
	pvzID := "pvz-1"

	mock.ExpectQuery(`(?i)SELECT p.id FROM products`).
		WithArgs(pvzID).
		WillReturnError(sql.ErrNoRows)

	err = services.DeleteLastProduct(db.DB, pvzID)
	require.Error(t, err)
	require.EqualError(t, err, "no products to delete or no active receptions")
	require.NoError(t, mock.ExpectationsWereMet())
}
