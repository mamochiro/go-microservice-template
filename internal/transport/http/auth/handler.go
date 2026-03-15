package auth

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/dto"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/handler"
)

type Handler struct {
	svc      service.AuthService
	validate *validator.Validate
}

func NewHandler(svc service.AuthService) *Handler {
	return &Handler{
		svc:      svc,
		validate: validator.New(),
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := handler.DecodeAndValidate(r, h.validate, &req); err != nil {
		handler.RespondError(w, err)
		return
	}

	accessToken, refreshToken, user, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		handler.RespondError(w, err)
		return
	}

	handler.RespondJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         handler.ToUserResponse(user),
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := handler.DecodeAndValidate(r, h.validate, &req); err != nil {
		handler.RespondError(w, err)
		return
	}

	accessToken, refreshToken, user, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		handler.RespondError(w, err)
		return
	}

	handler.RespondJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         handler.ToUserResponse(user),
	})
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := handler.DecodeAndValidate(r, h.validate, &req); err != nil {
		handler.RespondError(w, err)
		return
	}

	if err := h.svc.ForgotPassword(r.Context(), req.Email); err != nil {
		handler.RespondError(w, err)
		return
	}

	handler.RespondJSON(w, http.StatusOK, map[string]string{"message": "If the email is registered, a password reset link will be sent."})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := handler.DecodeAndValidate(r, h.validate, &req); err != nil {
		handler.RespondError(w, err)
		return
	}

	if err := h.svc.ResetPassword(r.Context(), req.Token, req.Password); err != nil {
		handler.RespondError(w, err)
		return
	}

	handler.RespondJSON(w, http.StatusOK, map[string]string{"message": "Password has been reset successfully."})
}
