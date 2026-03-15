package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/dto"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/handler"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/middleware"
)

const errInvalidID = "invalid id"

type Handler struct {
	svc      service.UserService
	validate *validator.Validate
}

func NewHandler(svc service.UserService) *Handler {
	v := validator.New()
	_ = v.RegisterValidation("nospaces", func(fl validator.FieldLevel) bool {
		return !strings.Contains(fl.Field().String(), " ")
	})

	return &Handler{
		svc:      svc,
		validate: v,
	}
}

func Register(r chi.Router, h *Handler, jwtSecret string) {
	// Public routes
	r.Post("/signup", h.Create)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))
		r.Route("/users", func(r chi.Router) {
			// Only Admins can list or delete users
			r.With(middleware.HasRole(entity.RoleAdmin)).Get("/", h.List)
			r.With(middleware.HasRole(entity.RoleAdmin)).Delete("/{id}", h.Delete)

			// Both User and Admin can get or update (further checks for "self" could be added)
			r.Get("/{id}", h.Get)
			r.Put("/{id}", h.Update)
		})
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := handler.DecodeAndValidate(r, h.validate, &req); err != nil {
		handler.RespondError(w, err)
		return
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.svc.CreateUser(r.Context(), user); err != nil {
		handler.RespondError(w, err)
		return
	}

	handler.RespondJSON(w, http.StatusCreated, handler.ToUserResponse(user))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	u, err := h.svc.GetUser(r.Context(), uint(id))
	if err != nil {
		handler.RespondError(w, err)
		return
	}

	handler.RespondJSON(w, http.StatusOK, handler.ToUserResponse(u))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	users, total, err := h.svc.ListUsersPaginated(r.Context(), page, limit)
	if err != nil {
		handler.RespondError(w, err)
		return
	}

	userRests := make([]dto.UserResponse, len(users))
	for i, u := range users {
		userRests[i] = handler.ToUserResponse(&u)
	}

	if limit < 1 {
		limit = 10
	}
	totalPages := int(total / int64(limit))
	if total%int64(limit) != 0 {
		totalPages++
	}

	handler.RespondJSON(w, http.StatusOK, dto.PaginatedUserResponse{
		Data:       userRests,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	var req dto.UpdateUserRequest
	if err := handler.DecodeAndValidate(r, h.validate, &req); err != nil {
		handler.RespondError(w, err)
		return
	}

	u := &entity.User{
		ID:       uint(id),
		Username: req.Username,
		Email:    req.Email,
	}

	if err := h.svc.UpdateUser(r.Context(), u); err != nil {
		handler.RespondError(w, err)
		return
	}

	handler.RespondJSON(w, http.StatusOK, handler.ToUserResponse(u))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteUser(r.Context(), uint(id)); err != nil {
		handler.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
