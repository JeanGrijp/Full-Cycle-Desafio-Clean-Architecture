package domain

import "time"

type Order struct {
	ID           int64
	CustomerName string
	Amount       float64
	Status       string
	CreatedAt    time.Time
}
