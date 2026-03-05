package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, func(), error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, nil, err
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
