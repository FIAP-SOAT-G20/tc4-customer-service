package entity

import (
	"time"
)

type Customer struct {
	ID        int
	Name      string
	Email     string
	CPF       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Customer) Update(name string, email string) {
	p.Name = name
	p.Email = email
	p.UpdatedAt = time.Now()
}
