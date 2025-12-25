package inventoryclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventorypb "github.com/DipeshC1707/BookYourShow/proto/inventory/v1"
)

type Client interface {
	LockSeats(
		ctx context.Context,
		eventID string,
		seatIDs []string,
		ownerID string,
	) error

	ReleaseSeats(
		ctx context.Context,
		eventID string,
		seatIDs []string,
	) error

	Close() error
}


type grpcClient struct {
	conn   *grpc.ClientConn
	client inventorypb.InventoryServiceClient
}

func New(address string) (Client, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		conn:   conn,
		client: inventorypb.NewInventoryServiceClient(conn),
	}, nil
}

func (c *grpcClient) LockSeats(
	ctx context.Context,
	eventID string,
	seatIDs []string,
	ownerID string,
) error {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.LockSeats(ctx, &inventorypb.LockSeatsRequest{
		EventId: eventID,
		SeatIds: seatIDs,
		OwnerId: ownerID,
	})
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf(resp.Error)
	}

	return nil
}

func (c *grpcClient) Close() error {
	return c.conn.Close()
}

func (c *grpcClient) ReleaseSeats(
	ctx context.Context,
	eventID string,
	seatIDs []string,
) error {

	_, err := c.client.ReleaseSeats(ctx, &inventorypb.ReleaseSeatsRequest{
		EventId: eventID,
		SeatIds: seatIDs,
	})

	return err
}