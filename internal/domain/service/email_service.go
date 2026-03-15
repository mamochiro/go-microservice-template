package service

import "context"

type EmailService interface {
	SendPasswordResetEmail(ctx context.Context, to, token string) error
}
