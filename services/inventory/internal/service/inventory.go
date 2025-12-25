package service

import (
	"context"
	"fmt"
	"time"

	redisclient "github.com/DipeshC1707/BookYourShow/inventory/internal/redis"
)

type InventoryService struct {
	redis *redisclient.Client
	lockTTL time.Duration
}

func NewInventoryService(
	redisClient *redisclient.Client,
	lockTTL time.Duration,
) *InventoryService {
	return &InventoryService{
		redis:   redisClient,
		lockTTL: lockTTL,
	}
}

func (s *InventoryService) LockSeats(
	ctx context.Context,
	eventID string,
	seatIDs []string,
	ownerID string,
) error {

	if len(seatIDs) == 0 {
		return fmt.Errorf("no seats provided")
	}

	keys := make([]string, 0, len(seatIDs))
	for _, seatID := range seatIDs {
		key := fmt.Sprintf("seat:%s:%s", eventID, seatID)
		keys = append(keys, key)
	}

	ok, err := s.redis.LockSeats(
		ctx,
		keys,
		ownerID,
		int(s.lockTTL.Seconds()),
	)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("one or more seats already locked")
	}

	return nil
}
