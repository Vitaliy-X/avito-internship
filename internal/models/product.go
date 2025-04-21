package models

import "time"

type Product struct {
	ID          string
	DateTime    time.Time
	Type        string
	ReceptionID string
}
