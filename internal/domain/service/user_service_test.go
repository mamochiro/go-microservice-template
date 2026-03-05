package service

import (
	"context"
	"testing"

	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := NewUserService(mockRepo)

	user := &entity.User{Username: "testuser", Email: "test@example.com"}

	// Setup expectation
	mockRepo.On("Create", mock.Anything, user).Return(nil)

	// Execute
	err := svc.CreateUser(context.Background(), user)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := NewUserService(mockRepo)

	expectedUser := &entity.User{ID: 1, Username: "testuser"}

	// Setup expectation
	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedUser, nil)

	// Execute
	user, err := svc.GetUser(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}
