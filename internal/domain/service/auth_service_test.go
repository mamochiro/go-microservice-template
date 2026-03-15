package service

import (
	"context"
	"testing"
	"time"

	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_ForgotPassword(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	mockCacheRepo := mocks.NewMockCacheRepository(t)
	mockEmailSvc := mocks.NewMockEmailService(t)
	cfg := &config.Config{App: config.AppConfig{JWTSecret: "secret"}}

	svc := NewAuthService(mockRepo, mockCacheRepo, mockEmailSvc, cfg)

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		user := &entity.User{ID: 1, Email: email}

		mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
		mockCacheRepo.On("Set", mock.Anything, mock.MatchedBy(func(key string) bool {
			return len(key) > 0
		}), email, 15*time.Minute).Return(nil)
		mockEmailSvc.On("SendPasswordResetEmail", mock.Anything, email, mock.Anything).Return(nil)

		err := svc.ForgotPassword(context.Background(), email)

		assert.NoError(t, err)
	})

	t.Run("user not found - should not return error", func(t *testing.T) {
		email := "nonexistent@example.com"
		mockRepo.On("GetByEmail", mock.Anything, email).Return(nil, assert.AnError)

		err := svc.ForgotPassword(context.Background(), email)

		assert.NoError(t, err)
	})
}

func TestAuthService_ResetPassword(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(t)
	mockCacheRepo := mocks.NewMockCacheRepository(t)
	mockEmailSvc := mocks.NewMockEmailService(t)
	cfg := &config.Config{App: config.AppConfig{JWTSecret: "secret"}}

	svc := NewAuthService(mockRepo, mockCacheRepo, mockEmailSvc, cfg)

	t.Run("success", func(t *testing.T) {
		token := "valid-token"
		email := "test@example.com"
		newPassword := "new-secure-password"
		user := &entity.User{ID: 1, Email: email, Password: "old-password"}

		mockCacheRepo.On("Get", mock.Anything, "reset_token:"+token).Return(email, nil)
		mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
		mockCacheRepo.On("Delete", mock.Anything, "reset_token:"+token).Return(nil)

		err := svc.ResetPassword(context.Background(), token, newPassword)

		assert.NoError(t, err)
		// Verify password was hashed and updated
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
		assert.NoError(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		token := "invalid-token"
		mockCacheRepo.On("Get", mock.Anything, "reset_token:"+token).Return("", assert.AnError)

		err := svc.ResetPassword(context.Background(), token, "password")

		assert.Error(t, err)
	})
}
