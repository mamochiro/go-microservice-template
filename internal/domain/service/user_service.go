package service

import (
	"context"

	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/domain/repository"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, id uint) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context) ([]entity.User, error)
	ListUsersPaginated(ctx context.Context, page, limit int) ([]entity.User, int64, error)
}

var tracer = otel.Tracer("user-service")

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) error {
	ctx, span := tracer.Start(ctx, "UserService.CreateUser")
	defer span.End()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.repo.Create(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	ctx, span := tracer.Start(ctx, "UserService.GetUser")
	defer span.End()
	return s.repo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	ctx, span := tracer.Start(ctx, "UserService.UpdateUser")
	defer span.End()
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	ctx, span := tracer.Start(ctx, "UserService.DeleteUser")
	defer span.End()
	return s.repo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]entity.User, error) {
	ctx, span := tracer.Start(ctx, "UserService.ListUsers")
	defer span.End()
	return s.repo.List(ctx)
}

func (s *userService) ListUsersPaginated(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	ctx, span := tracer.Start(ctx, "UserService.ListUsersPaginated")
	defer span.End()

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.ListPaginated(ctx, offset, limit)
}
