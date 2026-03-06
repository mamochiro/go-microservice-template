package repository

import (
	"context"

	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

var tracer = otel.Tracer("user-repository")

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	ctx, span := tracer.Start(ctx, "UserRepository.Create")
	defer span.End()
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	ctx, span := tracer.Start(ctx, "UserRepository.GetByID")
	defer span.End()

	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctx, span := tracer.Start(ctx, "UserRepository.GetByEmail")
	defer span.End()

	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	ctx, span := tracer.Start(ctx, "UserRepository.Update")
	defer span.End()
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	ctx, span := tracer.Start(ctx, "UserRepository.Delete")
	defer span.End()
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *UserRepo) List(ctx context.Context) ([]entity.User, error) {
	ctx, span := tracer.Start(ctx, "UserRepository.List")
	defer span.End()

	var users []entity.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) ListPaginated(ctx context.Context, offset, limit int) ([]entity.User, int64, error) {
	ctx, span := tracer.Start(ctx, "UserRepository.ListPaginated")
	defer span.End()

	var users []entity.User
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.User{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
