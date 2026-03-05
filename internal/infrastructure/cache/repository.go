package cache

import (
	"context"
	"time"

	"github.com/mamochiro/go-microservice-template/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) repository.CacheRepository {
	return &redisRepo{client: client}
}

func (r *redisRepo) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisRepo) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
