package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, func(), error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	cleanup := func() {
		err := client.Close()
		if err != nil {
			return
		}
		log.Println("Redis connection closed")
	}

	log.Println("Connected to Redis")
	return client, cleanup, nil
}
