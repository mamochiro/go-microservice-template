package service

import (
	"context"

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

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.List(ctx)
}
