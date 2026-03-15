package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/dto"
	"github.com/mamochiro/go-microservice-template/pkg/apperror"
)

// RespondJSON writes a JSON response with the given status code.
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// RespondError handles error mapping and writes a JSON error response.
func RespondError(w http.ResponseWriter, err error) {
	// Simple error mapping for now. Can be expanded.
	if err.Error() == "record not found" { // GORM specific
		RespondJSON(w, http.StatusNotFound, map[string]string{"error": "resource not found"})
		return
	}

	// Handle validator.ValidationErrors
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msg := formatValidationError(ve)
		RespondJSON(w, apperror.ErrValidation.Code, map[string]string{"error": msg})
		return
	}

	appErr := apperror.Convert(err)
	RespondJSON(w, appErr.Code, map[string]string{"error": appErr.Message})
}

// formatValidationError converts validator errors into human-readable messages.
func formatValidationError(ve validator.ValidationErrors) string {
	err := ve[0] // Take the first error
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "nospaces":
		return fmt.Sprintf("%s cannot contain spaces", err.Field())
	default:
		return fmt.Sprintf("Invalid value for %s", err.Field())
	}
}

// DecodeAndValidate decodes the request body and validates the struct.
func DecodeAndValidate(r *http.Request, v *validator.Validate, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperror.ErrBadRequest
	}
	if err := v.Struct(dst); err != nil {
		return err
	}
	return nil
}

// ToUserResponse maps a user entity to a user response DTO.
func ToUserResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
