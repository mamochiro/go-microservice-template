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
	"testing"

	"github.com/mamochiro/go-microservice-template/internal/app"
	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"

	"github.com/mamochiro/go-microservice-template/internal/transport/http/dto"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	server  *httptest.Server
	cleanup func()
}

func (s *APITestSuite) SetupSuite() {
	originalDir, err := os.Getwd()
	s.Require().NoError(err)

	err = os.Chdir("../..")
	s.Require().NoError(err)

	cfg, err := config.LoadConfig()
	s.Require().NoError(err)

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
	baseURL := s.server.URL + "/api/v1/users"

	// Create User
	user := &entity.User{
		Username: "api_test_user",
		Email:    "api_test@example.com",
	}
	body, _ := json.Marshal(user)

	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(body))
	s.NoError(err)
	s.Equal(http.StatusCreated, resp.StatusCode)

	var createdUser entity.User
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	s.NoError(err)
	s.NotZero(createdUser.ID)
	s.Equal(user.Username, createdUser.Username)

	// Get User
	resp, err = http.Get(fmt.Sprintf("%s/%d", baseURL, createdUser.ID))
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var foundUser entity.User
	err = json.NewDecoder(resp.Body).Decode(&foundUser)
	s.NoError(err)
	s.Equal(createdUser.ID, foundUser.ID)
	s.Equal(createdUser.Username, foundUser.Username)

	// Update User
	foundUser.Username = "updated_api_user"
	body, _ = json.Marshal(foundUser)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", baseURL, foundUser.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	// List Users
	resp, err = http.Get(baseURL)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var paginatedResp dto.PaginatedUserResponse
	err = json.NewDecoder(resp.Body).Decode(&paginatedResp)
	s.NoError(err)
	s.GreaterOrEqual(paginatedResp.Total, int64(1))
	s.NotEmpty(paginatedResp.Data)
	// Delete User
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", baseURL, foundUser.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	s.NoError(err)
	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
