package entity

import (
	"time"
)

type Customer struct {
	ID        uint64
	Name      string
	Email     string
	CPF       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
