package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/go-playground/validator/v10"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/dto"
	"github.com/mamochiro/go-microservice-template/pkg/apperror"
)

const errInvalidID = "invalid id"

type UserHandler struct {
	svc      service.UserService
	validate *validator.Validate
}

func NewUserHandler(svc service.UserService) *UserHandler {
	v := validator.New()

	// Register custom validation
	err := v.RegisterValidation("nospaces", func(fl validator.FieldLevel) bool {
		return !strings.Contains(fl.Field().String(), " ")
	})
	if err != nil {
		return nil
	}

	return &UserHandler{
		svc:      svc,
		validate: v,
	}
}

// Create handles the user creation request.
// @Summary      Create a new user (Signup)
// @Description  Create a new user with the provided details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserRequest  true  "User object"
// @Success      201   {object}  dto.UserResponse
// @Failure      400   {string}  string "Invalid request body"
// @Failure      422   {string}  string "Validation failed"
// @Failure      500   {string}  string "Internal server error"
// @Router       /signup [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, apperror.ErrBadRequest.Message, apperror.ErrBadRequest.Code)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), apperror.ErrValidation.Code)
		return
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // In a real app, hash this here or in service
	}

	if err := h.svc.CreateUser(r.Context(), user); err != nil {
		respondError(w, err)
		return
	}

	resp := toUserResponse(user)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}

// Get handles the request to retrieve a user by ID.
// @Summary      Get a user
// @Description  Get a user's details by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {string}  string "Invalid ID"
// @Failure      404  {string}  string "User not found"
// @Router       /users/{id} [get]
// @Security     Bearer
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetUser(r.Context(), uint(id))
	if err != nil {
		respondError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(toUserResponse(user)); err != nil {
		return
	}
}

// List handles the request to retrieve users with pagination.
// @Summary      List users
// @Description  Get a list of users with pagination
// @Tags         users
// @Produce      json
// @Param        page   query     int  false  "Page number (default 1)"
// @Param        limit  query     int  false  "Items per page (default 10)"
// @Success      200    {object}  dto.PaginatedUserResponse
// @Failure      500    {string}  string "Internal server error"
// @Security     Bearer
// @Router       /users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	users, total, err := h.svc.ListUsersPaginated(r.Context(), page, limit)
	if err != nil {
		respondError(w, err)
		return
	}

	userResps := make([]dto.UserResponse, len(users))
	for i, u := range users {
		userResps[i] = toUserResponse(&u)
	}

	// Calculate total pages
	if limit < 1 {
		limit = 10
	}
	totalPages := int(total / int64(limit))
	if total%int64(limit) != 0 {
		totalPages++
	}

	resp := dto.PaginatedUserResponse{
		Data:       userResps,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}

// @Summary      Update a user
// @Description  Update a user's details by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "User ID"
// @Param        user  body      dto.UpdateUserRequest  true  "Updated user object"
// @Success      200   {object}  dto.UserResponse
// @Failure      400   {string}  string "Invalid ID"
// @Failure      422   {string}  string "Validation failed"
// @Failure      500   {string}  string "Internal server error"
// @Security     Bearer
// @Router       /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, apperror.ErrBadRequest.Message, apperror.ErrBadRequest.Code)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), apperror.ErrValidation.Code)
		return
	}

	// First get existing to update only fields present?
	// Or just pass what we have. Service handles logic.
	// For simplicity, we map DTO to entity.
	user := &entity.User{
		ID:       uint(id),
		Username: req.Username,
		Email:    req.Email,
	}

	if err := h.svc.UpdateUser(r.Context(), user); err != nil {
		respondError(w, err)
		return
	}

	// Fetch updated user to return full object? Or just return what we have.
	// Let's assume UpdateUser updates the passed pointer or we can fetch again.
	// For now, return what we have.
	if err := json.NewEncoder(w).Encode(toUserResponse(user)); err != nil {
		return
	}
}

// @Summary      Delete a user
// @Description  Delete a user's record by their ID
// @Tags         users
// @Param        id   path      int  true  "User ID"
// @Success      204  {string}  string "No Content"
// @Failure      400  {string}  string "Invalid ID"
// @Failure      500  {string}  string "Internal server error"
// @Security     Bearer
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteUser(r.Context(), uint(id)); err != nil {
		respondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toUserResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func respondError(w http.ResponseWriter, err error) {
	// Simple error mapping for now. Can be expanded.
	// If it's a specific AppError, use its code.
	// Otherwise 500 or 404 depending on error.
	if err.Error() == "record not found" { // GORM specific
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	appErr := apperror.Convert(err)
	http.Error(w, appErr.Message, appErr.Code)
}
