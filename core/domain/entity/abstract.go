package entity

import "time"

type AbstractEntity struct {
	ID string

	CreatedAt time.Time
	UpdatedAt time.Time
}
