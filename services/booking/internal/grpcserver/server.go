package grpcserver

import (
	"context"

	bookingpb "github.com/DipeshC1707/BookYourShow/proto/booking/v1"
	"github.com/DipeshC1707/BookYourShow/booking/internal/service"
)

type Server struct {
	bookingpb.UnimplementedBookingServiceServer
	bookingService *service.BookingService
}

func New(bookingService *service.BookingService) *Server {
	return &Server{
		bookingService: bookingService,
	}
}

func (s *Server) CreateBooking(
	ctx context.Context,
	req *bookingpb.CreateBookingRequest,
) (*bookingpb.CreateBookingResponse, error) {

	bookingID, err := s.bookingService.CreateBooking(
		ctx,
		req.EventId,
		req.SeatIds,
		req.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &bookingpb.CreateBookingResponse{
		BookingId: bookingID,
	}, nil
}

func (s *Server) ConfirmBooking(
	ctx context.Context,
	req *bookingpb.ConfirmBookingRequest,
) (*bookingpb.ConfirmBookingResponse, error) {

	if err := s.bookingService.ConfirmBooking(ctx, req.BookingId); err != nil {
		return nil, err
	}

	return &bookingpb.ConfirmBookingResponse{
		Success: true,
	}, nil
}
