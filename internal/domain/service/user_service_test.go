package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	mockCache := mocks.NewMockCacheRepository(t)
	svc := NewUserService(mockRepo, mockCache)

	user := &entity.User{Username: "tester", Email: "test@example.com"}

	// Setup expectation
	mockRepo.On("Create", mock.Anything, user).Return(nil)

	// Execute
	err := svc.CreateUser(context.Background(), user)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_GetUser_CacheMiss(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	mockCache := mocks.NewMockCacheRepository(t)
	svc := NewUserService(mockRepo, mockCache)

	expectedUser := &entity.User{ID: 1, Username: "testuser"}
	cacheKey := fmt.Sprintf(userCacheKeyFormat, 1)

	// 1. Cache Get - Miss
	mockCache.On("Get", mock.Anything, cacheKey).Return("", fmt.Errorf("not found"))

	// 2. Repo GetByID - Success
	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedUser, nil)

	// 3. Cache Set - Success
	userJSON, _ := json.Marshal(expectedUser)
	mockCache.On("Set", mock.Anything, cacheKey, userJSON, 10*time.Minute).Return(nil)

	// Execute
	user, err := svc.GetUser(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_GetUser_CacheHit(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	mockCache := mocks.NewMockCacheRepository(t)
	svc := NewUserService(mockRepo, mockCache)

	expectedUser := &entity.User{ID: 1, Username: "testuser"}
	cacheKey := fmt.Sprintf(userCacheKeyFormat, 1)
	userJSON, _ := json.Marshal(expectedUser)

	// 1. Cache Get - Hit
	mockCache.On("Get", mock.Anything, cacheKey).Return(string(userJSON), nil)

	// Execute
	user, err := svc.GetUser(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	// Repo should NOT be called
	mockRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestUserService_ListUsersPaginated(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	mockCache := mocks.NewMockCacheRepository(t)
	svc := NewUserService(mockRepo, mockCache)

	expectedUsers := []entity.User{{ID: 1, Username: "user1"}}

	mockRepo.On("ListPaginated", mock.Anything, 0, 10).Return(expectedUsers, int64(1), nil)

	users, total, err := svc.ListUsersPaginated(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, users, 1)
}
