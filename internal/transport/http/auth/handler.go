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
