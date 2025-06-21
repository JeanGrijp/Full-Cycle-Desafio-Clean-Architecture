package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/graph"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/repository"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/usecase"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := sql.Open("postgres", "postgres://root:root@localhost:5432/ordersdb?sslmode=disable")
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}
	defer db.Close()

	orderRepo := &repository.OrderPgRepository{DB: db}

	orderUseCase := &usecase.ListOrdersUseCase{Repo: orderRepo}

	resolver := &graph.Resolver{
		OrderUseCase: orderUseCase,
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
