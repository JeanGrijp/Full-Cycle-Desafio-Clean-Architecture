package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/repository"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/usecase"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Pegando configs do banco via ENV
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Erro ao conectar ao banco de dados", "error", err)
		log.Fatal(err)
	}
	defer db.Close()

	orderRepo := &repository.OrderPgRepository{DB: db}
	orderUseCase := &usecase.ListOrdersUseCase{Repo: orderRepo}

	slog.InfoContext(ctx, "Iniciando o servidor HTTP")

	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.RedirectSlashes)
	r.Use(chiMiddleware.Timeout(190 * time.Second))

	r.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(ctx, "Listando pedidos")
		orders, err := orderUseCase.Execute()
		slog.InfoContext(ctx, "Pedidos listados com sucesso", "count", len(orders))
		if err != nil {
			slog.ErrorContext(ctx, "Erro ao listar pedidos", "error", err)
			http.Error(w, "Erro ao listar pedidos: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if len(orders) == 0 {
			slog.InfoContext(ctx, "Nenhum pedido encontrado")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})

	log.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
