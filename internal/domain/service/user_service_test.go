package service

import (
	"context"
	"testing"

	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	svc := NewUserService(mockRepo)

	user := &entity.User{Username: "tester", Email: "test@example.com"}

	// Setup expectation
	mockRepo.On("Create", mock.Anything, user).Return(nil)

	// Execute
	err := svc.CreateUser(context.Background(), user)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	svc := NewUserService(mockRepo)

	expectedUser := &entity.User{ID: 1, Username: "testuser"}

	// Setup expectation
	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedUser, nil)

	// Execute
	user, err := svc.GetUser(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_ListUsersPaginated(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	svc := NewUserService(mockRepo)

	expectedUsers := []entity.User{{ID: 1, Username: "user1"}}

	mockRepo.On("ListPaginated", mock.Anything, 0, 10).Return(expectedUsers, int64(1), nil)

	users, total, err := svc.ListUsersPaginated(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, users, 1)
}
