package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/domain/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	userCacheKeyFormat      = "user:%d"
	userEmailCacheKeyFormat = "user:email:%s"
)

type cachedUserRepo struct {
	dbRepo *UserRepo
	cache  repository.CacheRepository
	tracer trace.Tracer
}

func NewCachedUserRepository(dbRepo *UserRepo, cache repository.CacheRepository) repository.UserRepository {
	return &cachedUserRepo{
		dbRepo: dbRepo,
		cache:  cache,
		tracer: otel.Tracer("cached-user-repository"),
	}
}

func (r *cachedUserRepo) Create(ctx context.Context, user *entity.User) error {
	return r.dbRepo.Create(ctx, user)
}

func (r *cachedUserRepo) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "CachedUserRepository.GetByID")
	defer span.End()

	cacheKey := fmt.Sprintf(userCacheKeyFormat, id)

	// 1. Try to get from cache
	val, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		var user entity.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			span.AddEvent("cache hit")
			return &user, nil
		}
	}

	span.AddEvent("cache miss")

	// 2. Fallback to DB
	user, err := r.dbRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. Save to cache for next time (expires in 10 minutes)
	userJSON, _ := json.Marshal(user)
	_ = r.cache.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	return user, nil
}

func (r *cachedUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "CachedUserRepository.GetByEmail")
	defer span.End()

	cacheKey := fmt.Sprintf(userEmailCacheKeyFormat, email)

	// 1. Try to get from cache
	val, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		var user entity.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			span.AddEvent("cache hit")
			return &user, nil
		}
	}

	span.AddEvent("cache miss")

	// 2. Fallback to DB
	user, err := r.dbRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// 3. Save to cache for next time
	userJSON, _ := json.Marshal(user)
	_ = r.cache.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	return user, nil
}

func (r *cachedUserRepo) Update(ctx context.Context, user *entity.User) error {
	ctx, span := r.tracer.Start(ctx, "CachedUserRepository.Update")
	defer span.End()

	if err := r.dbRepo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, fmt.Sprintf(userCacheKeyFormat, user.ID))
	_ = r.cache.Delete(ctx, fmt.Sprintf(userEmailCacheKeyFormat, user.Email))
	return nil
}

func (r *cachedUserRepo) Delete(ctx context.Context, id uint) error {
	ctx, span := r.tracer.Start(ctx, "CachedUserRepository.Delete")
	defer span.End()

	// Need to get user to invalidate email cache
	user, err := r.dbRepo.GetByID(ctx, id)
	if err == nil {
		_ = r.cache.Delete(ctx, fmt.Sprintf(userEmailCacheKeyFormat, user.Email))
	}

	if err := r.dbRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, fmt.Sprintf(userCacheKeyFormat, id))
	return nil
}

func (r *cachedUserRepo) List(ctx context.Context) ([]entity.User, error) {
	return r.dbRepo.List(ctx)
}

func (r *cachedUserRepo) ListPaginated(ctx context.Context, offset, limit int) ([]entity.User, int64, error) {
	return r.dbRepo.ListPaginated(ctx, offset, limit)
}
