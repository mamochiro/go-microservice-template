package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/domain/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, id uint) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context) ([]entity.User, error)
}

const userCacheKeyFormat = "user:%d"

type userService struct {
	repo  repository.UserRepository
	cache repository.CacheRepository
}

func NewUserService(repo repository.UserRepository, cache repository.CacheRepository) UserService {
	return &userService{repo: repo, cache: cache}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	cacheKey := fmt.Sprintf(userCacheKeyFormat, id)

	// 1. Try to get from cache
	val, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var user entity.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil
		}
	}

	// 2. Fallback to DB
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. Save to cache for next time (expires in 10 minutes)
	userJSON, _ := json.Marshal(user)
	_ = s.cache.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf(userCacheKeyFormat, user.ID))
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf(userCacheKeyFormat, id))
	return nil
}

func (s *userService) ListUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.List(ctx)
}
