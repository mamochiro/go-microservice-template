package auth

import "github.com/go-chi/chi/v5"

func Register(r chi.Router, h *Handler) {
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
}
