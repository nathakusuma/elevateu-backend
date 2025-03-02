package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/redis/go-redis/v9"
)

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, val interface{}) error
	Del(ctx context.Context, key string) error
	Close() error
}

type redisImpl struct {
	client *redis.Client
}

func NewRedis(host, port, pass string, db int) ICache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(map[string]interface{}{
			"error": err.Error(),
		}, "[REDIS][NewRedis] failed to connect to redis")
	}

	return &redisImpl{
		client: client,
	}
}

func (r *redisImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisImpl) Get(ctx context.Context, key string, val interface{}) error {
	err := r.client.Get(ctx, key).Scan(val)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("not found: %w", err)
		}
		return err
	}

	return nil
}

func (r *redisImpl) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisImpl) Close() error {
	return r.client.Close()
}
