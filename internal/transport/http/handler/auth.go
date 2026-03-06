package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/dto"
	"github.com/mamochiro/go-microservice-template/pkg/apperror"
)

type AuthHandler struct {
	svc      service.AuthService
	validate *validator.Validate
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc:      svc,
		validate: validator.New(),
	}
}

// Login handles the user login request.
// @Summary      User login
// @Description  Authenticate user and return access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login  body      dto.LoginRequest  true  "Login credentials"
// @Success      200    {object}  dto.AuthResponse
// @Failure      401    {string}  string "Invalid credentials"
// @Router       /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, apperror.ErrBadRequest.Message, apperror.ErrBadRequest.Code)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), apperror.ErrValidation.Code)
		return
	}

	token, user, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, err)
		return
	}

	resp := dto.AuthResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Since we already sent 200 OK, we just log the error
		return
	}
}
