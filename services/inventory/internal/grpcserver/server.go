package grpcserver

import (
	"context"
	"github.com/DipeshC1707/BookYourShow/proto/inventory/v1"
	"github.com/DipeshC1707/BookYourShow/inventory/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	inventorypb.UnimplementedInventoryServiceServer
	inventoryService *service.InventoryService
}

func NewServer(inventoryService *service.InventoryService) *Server {
	return &Server{
		inventoryService: inventoryService,
	}
}

func (s *Server) LockSeats(
	ctx context.Context,
	req *inventorypb.LockSeatsRequest,
) (*inventorypb.LockSeatsResponse, error) {

	err := s.inventoryService.LockSeats(
		ctx,
		req.EventId,
		req.SeatIds,
		req.OwnerId,
	)

	if err != nil {
		return &inventorypb.LockSeatsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &inventorypb.LockSeatsResponse{
		Success: true,
	}, nil
}

func (s *Server) ReleaseSeats(
	ctx context.Context,
	req *inventorypb.ReleaseSeatsRequest,
) (*emptypb.Empty, error) {

	err := s.inventoryService.ReleaseSeats(ctx, req.EventId, req.SeatIds)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}