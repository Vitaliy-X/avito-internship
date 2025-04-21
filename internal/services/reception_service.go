package services

import (
	"avito-internship/internal/models"
	"database/sql"
	"errors"
)

func CreateReception(db *sql.DB, pvzID string) (*models.Reception, error) {
	var existing string
	err := db.QueryRow(`
        SELECT id FROM receptions
        WHERE pvz_id = $1 AND status = 'in_progress'
        LIMIT 1
    `, pvzID).Scan(&existing)
	if err == nil {
		return nil, errors.New("already an open reception")
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	var reception models.Reception
	row := db.QueryRow(`
        INSERT INTO receptions (pvz_id, status)
        VALUES ($1, 'in_progress')
        RETURNING id, date_time, status
    `, pvzID)
	if err := row.Scan(&reception.ID, &reception.DateTime, &reception.Status); err != nil {
		return nil, err
	}
	reception.PVZID = pvzID
	return &reception, nil
}

func CloseLastReception(db *sql.DB, pvzID string) error {
	res, err := db.Exec(`
        UPDATE receptions
        SET status = 'close'
        WHERE id = (
            SELECT id FROM receptions
            WHERE pvz_id = $1 AND status = 'in_progress'
            ORDER BY date_time DESC
            LIMIT 1
        )
    `, pvzID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("no active reception")
	}
	return nil
}
