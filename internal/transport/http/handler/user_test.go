package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Create(t *testing.T) {
	mockSvc := mocks.NewMockUserService(t)
	h := NewUserHandler(mockSvc)

	user := &entity.User{Username: "testuser", Email: "test@example.com"}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mockSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	h.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
