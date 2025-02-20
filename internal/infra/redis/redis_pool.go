package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

var (
	client *redis.Client
	once   sync.Once
)

func NewRedisPool(host, port, pass string, db int) *redis.Client {
	once.Do(func() {
		cl := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: pass,
			DB:       db,
		})

		_, err := cl.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal(map[string]interface{}{
				"error": err.Error(),
			}, "[REDIS][NewRedisPool] failed to connect to redis")
		}

		client = cl
	})

	return client
}
