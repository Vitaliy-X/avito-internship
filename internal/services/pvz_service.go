package services

import (
	"avito-internship/internal/models"
	"database/sql"
	"errors"
	"log"
	"time"
)

var allowedCities = map[string]bool{
	"Москва": true, "Санкт-Петербург": true, "Казань": true,
}

type PVZWithReceptions struct {
	PVZ        models.PVZ
	Receptions []ReceptionWithProducts
}

type ReceptionWithProducts struct {
	Reception models.Reception
	Products  []models.Product
}

func CreatePVZ(db *sql.DB, pvz models.PVZ) (*models.PVZ, error) {
	if !allowedCities[pvz.City] {
		return nil, errors.New("city not allowed")
	}

	query := `
		INSERT INTO pvz (id, registration_date, city)
		VALUES ($1, $2, $3)
		RETURNING id, registration_date, city
	`

	row := db.QueryRow(query, pvz.ID, pvz.RegistrationDate, pvz.City)

	var newPVZ models.PVZ
	if err := row.Scan(&newPVZ.ID, &newPVZ.RegistrationDate, &newPVZ.City); err != nil {
		return nil, err
	}

	return &newPVZ, nil
}

func GetPVZList(db *sql.DB, startDate, endDate *time.Time, page, limit int) ([]PVZWithReceptions, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, registration_date, city
		FROM pvz
		WHERE ($1::timestamptz IS NULL OR registration_date >= $1::timestamptz)
		  AND ($2::timestamptz IS NULL OR registration_date <= $2::timestamptz)
		ORDER BY registration_date DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PVZWithReceptions

	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, err
		}

		recs, err := getReceptionsWithProducts(db, pvz.ID, startDate, endDate)
		if err != nil {
			return nil, err
		}

		results = append(results, PVZWithReceptions{
			PVZ:        pvz,
			Receptions: recs,
		})
	}

	return results, nil
}

func getReceptionsWithProducts(db *sql.DB, pvzID string, startDate, endDate *time.Time) ([]ReceptionWithProducts, error) {
	query := `
        SELECT id, date_time, status
        FROM receptions
        WHERE pvz_id = $1
          AND ($2::timestamptz IS NULL OR date_time >= $2::timestamptz)
          AND ($3::timestamptz IS NULL OR date_time <= $3::timestamptz)
        ORDER BY date_time
    `
	rows, err := db.Query(query, pvzID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var receptions []ReceptionWithProducts

	for rows.Next() {
		var r models.Reception
		if err := rows.Scan(&r.ID, &r.DateTime, &r.Status); err != nil {
			return nil, err
		}
		r.PVZID = pvzID

		products, err := getProductsByReception(db, r.ID)
		if err != nil {
			return nil, err
		}

		receptions = append(receptions, ReceptionWithProducts{
			Reception: r,
			Products:  products,
		})
	}

	return receptions, nil
}

func getProductsByReception(db *sql.DB, receptionID string) ([]models.Product, error) {
	rows, err := db.Query(`
        SELECT id, date_time, type
        FROM products
        WHERE reception_id = $1
        ORDER BY date_time
    `, receptionID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.DateTime, &p.Type); err != nil {
			return nil, err
		}
		p.ReceptionID = receptionID
		products = append(products, p)
	}

	return products, nil
}
