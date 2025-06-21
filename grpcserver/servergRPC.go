// grpcserver/servergRPC.go
package grpcserver

import (
	"context"
	"time"

	orderpb "github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/infra/grpc/pb"
	"github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/usecase"
)

type OrderServiceServer struct {
	orderpb.UnimplementedOrderServiceServer
	OrderUseCase usecase.OrderUseCaseInterface
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	orders, err := s.OrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var pbOrders []*orderpb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &orderpb.Order{
			Id:           o.ID,
			CustomerName: o.CustomerName,
			Amount:       o.Amount,
			Status:       o.Status,
			CreatedAt:    o.CreatedAt.Format(time.RFC3339),
		})
	}

	return &orderpb.ListOrdersResponse{Orders: pbOrders}, nil
}
