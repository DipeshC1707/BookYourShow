package bookingclient

import (
	"context"
	"time"

	bookingpb "github.com/DipeshC1707/BookYourShow/proto/booking/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client bookingpb.BookingServiceClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: bookingpb.NewBookingServiceClient(conn),
	}, nil
}

func (c *Client) ConfirmBooking(
	ctx context.Context,
	bookingID string,
) error {
	_, err := c.client.ConfirmBooking(ctx, &bookingpb.ConfirmBookingRequest{
		BookingId: bookingID,
	})
	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
