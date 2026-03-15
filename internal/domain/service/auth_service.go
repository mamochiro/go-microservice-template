package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/domain/repository"
	"github.com/mamochiro/go-microservice-template/pkg/apperror"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenKeyFormat  = "refresh_token:%s"
	resetTokenKeyFormat    = "reset_token:%s"
	accessTokenDuration    = time.Minute * 15
	refreshTokenDuration   = time.Hour * 24 * 7
	resetTokenDuration     = time.Minute * 15
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, *entity.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, *entity.User, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type authService struct {
	repo      repository.UserRepository
	cacheRepo repository.CacheRepository
	emailSvc  EmailService
	jwtSecret string
}

func NewAuthService(
	repo repository.UserRepository,
	cacheRepo repository.CacheRepository,
	emailSvc EmailService,
	cfg *config.Config,
) AuthService {
	return &authService{
		repo:      repo,
		cacheRepo: cacheRepo,
		emailSvc:  emailSvc,
		jwtSecret: cfg.App.JWTSecret,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, *entity.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", nil, apperror.New("invalid email or password", 401)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, apperror.New("invalid email or password", 401)
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, string, *entity.User, error) {
	// 1. Verify Refresh Token from Cache
	userIDStr, err := s.cacheRepo.Get(ctx, fmt.Sprintf(refreshTokenKeyFormat, refreshToken))
	if err != nil {
		return "", "", nil, apperror.New("invalid or expired refresh token", 401)
	}

	// 2. Get User
	var userID uint
	_, err = fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to parse user id from cache: %w", err)
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return "", "", nil, apperror.New("user not found", 404)
	}

	// 3. Delete old refresh token (rotate)
	_ = s.cacheRepo.Delete(ctx, fmt.Sprintf(refreshTokenKeyFormat, refreshToken))

	// 4. Generate new tokens
	accessToken, newRefreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, newRefreshToken, user, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// We return nil to prevent email enumeration attacks
		return nil
	}

	// Generate a secure reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store token in cache (expires in 15 mins)
	err = s.cacheRepo.Set(ctx, fmt.Sprintf(resetTokenKeyFormat, token), user.Email, resetTokenDuration)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Send email
	return s.emailSvc.SendPasswordResetEmail(ctx, user.Email, token)
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// 1. Get email from token
	email, err := s.cacheRepo.Get(ctx, fmt.Sprintf(resetTokenKeyFormat, token))
	if err != nil {
		return apperror.New("invalid or expired reset token", 400)
	}

	// 2. Get user
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return apperror.New("user not found", 404)
	}

	// 3. Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// 4. Update user
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	// 5. Delete reset token
	_ = s.cacheRepo.Delete(ctx, fmt.Sprintf(resetTokenKeyFormat, token))

	return nil
}

func (s *authService) generateTokens(ctx context.Context, user *entity.User) (string, string, error) {
	// Generate Access Token (Short-lived)
	accessTokenClaims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(accessTokenDuration).Unix(),
		"iat":  time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate Refresh Token (Long-lived)
	refreshTokenClaims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(refreshTokenDuration).Unix(),
		"iat":  time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// Store Refresh Token in Cache
	err = s.cacheRepo.Set(ctx, fmt.Sprintf(refreshTokenKeyFormat, refreshTokenString), user.ID, refreshTokenDuration)
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token in cache: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}
