//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mamochiro/go-microservice-template/internal/app"
	"github.com/mamochiro/go-microservice-template/internal/config"

	"github.com/mamochiro/go-microservice-template/internal/transport/http/dto"
	"github.com/mamochiro/go-microservice-template/pkg/logger"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	server  *httptest.Server
	cleanup func()
}

func (s *APITestSuite) SetupSuite() {
	// Initialize logger for tests
	logger.Init("test")

	// Find project root by looking for go.mod
	originalDir, err := os.Getwd()
	s.Require().NoError(err)

	currentDir := originalDir
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			s.Fail("Could not find project root (go.mod)")
		}
		currentDir = parent
	}

	// Change to root to load config
	err = os.Chdir(currentDir)
	s.Require().NoError(err)

	cfg, err := config.LoadConfig()
	s.Require().NoError(err)

	// Set migration path to absolute path
	absMigrationPath, err := filepath.Abs("migrations")
	s.Require().NoError(err)
	cfg.Postgres.MigrationPath = absMigrationPath

	err = os.Chdir(originalDir)
	s.Require().NoError(err)

	mux, cleanup, err := app.InitializeApp(cfg)
	s.Require().NoError(err)

	s.server = httptest.NewServer(mux)
	s.cleanup = cleanup
}

func (s *APITestSuite) TearDownSuite() {
	s.server.Close()
	if s.cleanup != nil {
		s.cleanup()
	}
}

func (s *APITestSuite) TestHealthCheck() {
	resp, err := http.Get(s.server.URL + "/health")
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *APITestSuite) TestUserLifecycle() {
	baseURL := s.server.URL + "/api/v1"

	// 1. Create User (Public)
	userReq := map[string]string{
		"username": "api_test_user",
		"email":    "api_test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(userReq)

	resp, err := http.Post(baseURL+"/signup", "application/json", bytes.NewBuffer(body))
	s.NoError(err)
	s.Equal(http.StatusCreated, resp.StatusCode)

	var createdUser dto.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	s.NoError(err)
	s.NotZero(createdUser.ID)

	// 2. Login to get token
	loginReq := dto.LoginRequest{
		Email:    userReq["email"],
		Password: userReq["password"],
	}
	body, _ = json.Marshal(loginReq)
	resp, err = http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(body))
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var authResp dto.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	s.NoError(err)
	s.NotEmpty(authResp.AccessToken)
	token := authResp.AccessToken

	// Helper function for authenticated requests
	doAuthRequest := func(method, url string, body []byte) *http.Response {
		req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		s.NoError(err)
		return res
	}

	// 3. Get User (Protected)
	resp = doAuthRequest(http.MethodGet, fmt.Sprintf("%s/users/%d", baseURL, createdUser.ID), nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	var foundUser dto.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&foundUser)
	s.NoError(err)
	s.Equal(createdUser.ID, foundUser.ID)

	// 4. Update User (Protected)
	updateReq := dto.UpdateUserRequest{
		Username: "updated_api_user",
	}
	body, _ = json.Marshal(updateReq)
	resp = doAuthRequest(http.MethodPut, fmt.Sprintf("%s/users/%d", baseURL, foundUser.ID), body)
	s.Equal(http.StatusOK, resp.StatusCode)

	// 5. List Users (Protected)
	resp = doAuthRequest(http.MethodGet, baseURL+"/users", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	var paginatedResp dto.PaginatedUserResponse
	err = json.NewDecoder(resp.Body).Decode(&paginatedResp)
	s.NoError(err)
	s.GreaterOrEqual(paginatedResp.Total, int64(1))

	// 6. Delete User (Protected)
	resp = doAuthRequest(http.MethodDelete, fmt.Sprintf("%s/users/%d", baseURL, foundUser.ID), nil)
	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
