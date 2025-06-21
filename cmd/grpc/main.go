package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/grpcserver"
	orderpb "github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/infra/grpc/pb"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/repository"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/usecase"
)

func main() {
	db, err := sql.Open("postgres", "postgres://usuario:senha@localhost:5432/sua_base?sslmode=disable")
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}
	defer db.Close()

	orderRepo := &repository.OrderPgRepository{DB: db}
	orderUseCase := &usecase.ListOrdersUseCase{Repo: orderRepo}

	grpcServer := grpc.NewServer()
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
