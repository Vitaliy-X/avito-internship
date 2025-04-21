package models

import "time"

type Reception struct {
	ID       string
	DateTime time.Time
	Status   string
	PVZID    string
}
