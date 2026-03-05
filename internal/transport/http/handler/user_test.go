package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/dto"
	"github.com/mamochiro/go-microservice-template/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Create(t *testing.T) {
	mockSvc := mocks.NewMockUserService(t)
	h := NewUserHandler(mockSvc)

	tests := []struct {
		name           string
		requestBody    dto.CreateUserRequest
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Success",
			requestBody: dto.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Validation Error",
			requestBody: dto.CreateUserRequest{
				Username: "ab", // Too short
				Email:    "invalid-email",
				Password: "pass", // Too short
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusUnprocessableEntity, // apperror.ErrValidation.Code
		},
		{
			name: "Service Error",
			requestBody: dto.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			h.Create(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	mockSvc := mocks.NewMockUserService(t)
	h := NewUserHandler(mockSvc)

	tests := []struct {
		name           string
		userID         string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "Success",
			userID: "1",
			mockSetup: func() {
				mockSvc.On("GetUser", mock.Anything, uint(1)).Return(&entity.User{ID: 1, Username: "user1"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid ID",
			userID: "abc",
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Not Found",
			userID: "99",
			mockSetup: func() {
				mockSvc.On("GetUser", mock.Anything, uint(99)).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError, // Default error mapping
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)

			// Inject chi context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			h.Get(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
