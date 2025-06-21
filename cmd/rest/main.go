package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/repository"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/usecase"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://root:root@localhost:5432/ordersdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	orderRepo := &repository.OrderPgRepository{DB: db}
	orderUseCase := &usecase.ListOrdersUseCase{Repo: orderRepo}

	r := chi.NewRouter()

	r.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
		orders, err := orderUseCase.Execute()
		if err != nil {
			http.Error(w, "Erro ao listar pedidos: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})

	log.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
