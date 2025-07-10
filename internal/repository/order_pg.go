package repository

import (
	"context"
	"database/sql"

	"log/slog"

	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/domain"
)

type OrderPgRepository struct {
	DB *sql.DB
}

func (r *OrderPgRepository) List() ([]domain.Order, error) {
	ctx := context.Background()
	rows, err := r.DB.QueryContext(ctx, "SELECT id, customer_name, amount, status, created_at FROM orders")
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao executar consulta", "error", err)
		return nil, err
	}

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Amount, &o.Status, &o.CreatedAt); err != nil {
			slog.ErrorContext(ctx, "Erro ao escanear pedido", "error", err)
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Close(); err != nil {
		slog.ErrorContext(ctx, "Erro ao fechar rows", "error", err)
	}
	slog.InfoContext(ctx, "Pedidos listados com sucesso", "count", len(orders))
	return orders, nil
}
