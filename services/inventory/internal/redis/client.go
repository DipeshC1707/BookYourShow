package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(redisAddr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	return &Client{rdb: rdb}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) LockSeats(
	ctx context.Context,
	keys []string,
	owner string,
	ttlSeconds int,
) (bool, error) {

	result, err := c.rdb.Eval(
		ctx,
		LockSeatsLua,
		keys,
		owner,
		ttlSeconds,
	).Int()

	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (c *Client) DeleteKeys(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.rdb.Del(ctx, keys...).Err()
}
