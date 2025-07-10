package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/grpcserver"
	orderpb "github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/infra/grpc/pb"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/repository"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/usecase"
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
		slog.ErrorContext(ctx, "Erro ao conectar ao banco de dados", "error", err)
		log.Fatal(err)
	}
	defer db.Close()

	orderRepo := &repository.OrderPgRepository{DB: db}
	orderUseCase := &usecase.ListOrdersUseCase{Repo: orderRepo}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	orderpb.RegisterOrderServiceServer(grpcServer, &grpcserver.OrderServiceServer{
		OrderUseCase: orderUseCase,
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao escutar: %v", err)
	}

	log.Println("gRPC server rodando na porta 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}
