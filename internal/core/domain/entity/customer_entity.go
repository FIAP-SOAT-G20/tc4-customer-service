package entity

import (
	"time"
)

type Customer struct {
	ID        string
	Name      string
	Email     string
	CPF       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
