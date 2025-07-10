package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/graph"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/repository"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/usecase"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

const defaultPort = "8081"

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Nova requisição recebida",
			"method", r.Method,
			"url", r.URL.Path,
			"content-type", r.Header.Get("Content-Type"),
		)
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Requisição finalizada",
			"method", r.Method,
			"url", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Pegando configs do banco via ENV
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		slog.Error("Alguma variável de ambiente do banco está vazia!",
			"user", dbUser, "password", dbPassword, "host", dbHost, "port", dbPort, "db", dbName)
		log.Fatal("Configuração do banco incompleta")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao conectar ao banco de dados", "error", err)
		log.Fatal(err)
	}
	defer db.Close()

	slog.InfoContext(ctx, "Banco de dados conectado com sucesso")

	orderRepo := &repository.OrderPgRepository{DB: db}
	orderUseCase := &usecase.ListOrdersUseCase{Repo: orderRepo}
	resolver := &graph.Resolver{OrderUseCase: orderUseCase}

	slog.InfoContext(ctx, "Iniciando o servidor GraphQL")

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// -- Recursos extras para produção e compatibilidade com todos os clients
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	// Chi router
	r := chi.NewRouter()
	r.Use(logMiddleware)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.RedirectSlashes)
	r.Use(chiMiddleware.Timeout(190 * time.Second))

	// Playground agora aponta para o endpoint correto "/query"
	r.Get("/playground", func(w http.ResponseWriter, r *http.Request) {
		playground.Handler("GraphQL playground", "/query").ServeHTTP(w, r)
	})
	r.Handle("/query", srv) // endpoint principal GraphQL

	slog.InfoContext(ctx, "Servidor HTTP rodando na porta", "port", port)
	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
