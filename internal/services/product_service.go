package services

import (
	"avito-internship/internal/models"
	"database/sql"
	"errors"
)

func AddProduct(db *sql.DB, pvzID, productType string) (*models.Product, error) {
	var receptionID string
	err := db.QueryRow(`
        SELECT id FROM receptions
        WHERE pvz_id = $1 AND status = 'in_progress'
        ORDER BY date_time DESC
        LIMIT 1
    `, pvzID).Scan(&receptionID)

	if err != nil {
		return nil, errors.New("no active product")
	}

	var product models.Product
	row := db.QueryRow(`
        INSERT INTO products (type, reception_id)
        VALUES ($1, $2)
        RETURNING id, date_time
    `, productType, receptionID)

	if err := row.Scan(&product.ID, &product.DateTime); err != nil {
		return nil, err
	}

	product.Type = productType
	product.ReceptionID = receptionID
	return &product, nil
}

func DeleteLastProduct(db *sql.DB, pvzID string) error {
	var productID string
	err := db.QueryRow(`
        SELECT p.id FROM products p
        JOIN receptions r ON r.id = p.reception_id
        WHERE r.pvz_id = $1 AND r.status = 'in_progress'
        ORDER BY p.date_time DESC
        LIMIT 1
    `, pvzID).Scan(&productID)

	if err != nil {
		return errors.New("no products to delete or no active receptions")
	}

	_, err = db.Exec(`DELETE FROM products WHERE id = $1`, productID)
	return err
}
