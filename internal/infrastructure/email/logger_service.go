package email

import (
	"context"
	"fmt"

	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/mamochiro/go-microservice-template/pkg/logger"
)

type loggerService struct{}

func NewLoggerService() service.EmailService {
	return &loggerService{}
}

func (s *loggerService) SendPasswordResetEmail(ctx context.Context, to, token string) error {
	logger.Info(fmt.Sprintf("Sending password reset email to %s with token: %s", to, token))
	return nil
}
