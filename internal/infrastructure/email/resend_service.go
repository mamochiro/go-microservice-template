package email

import (
	"context"
	"fmt"

	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/resend/resend-go/v2"
)

type resendService struct {
	client *resend.Client
	from   string
}

func NewResendService(cfg *config.Config) service.EmailService {
	client := resend.NewClient(cfg.Email.ApiKey)
	return &resendService{
		client: client,
		from:   cfg.Email.From,
	}
}

func (s *resendService) SendPasswordResetEmail(ctx context.Context, to, token string) error {
	// In a real app, this would be a URL to your frontend
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: "Password Reset Request",
		Html:    fmt.Sprintf("<strong>Reset your password:</strong> <a href='%s'>Click here</a>", resetURL),
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	return err
}
