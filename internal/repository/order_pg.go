package repository

import (
	"database/sql"

	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/domain"
)

type OrderPgRepository struct {
	DB *sql.DB
}

func (r *OrderPgRepository) List() ([]domain.Order, error) {
	rows, err := r.DB.Query("SELECT id, customer_name, amount, status, created_at FROM orders")
	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Amount, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
