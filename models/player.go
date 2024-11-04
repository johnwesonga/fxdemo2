package models

import "time"

type Player struct {
	ID        int       `json:"Id"`
	Name      string    `json:"Name"`
	Score     int       `json:"Score"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
