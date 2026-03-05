//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/database"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db      *gorm.DB
	cleanup func()
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	// Temporarily change directory to root to load config
	originalDir, err := os.Getwd()
	s.Require().NoError(err)

	err = os.Chdir("../..")
	s.Require().NoError(err)

	cfg, err := config.LoadConfig()
	s.Require().NoError(err)

	// Change back
	err = os.Chdir(originalDir)
	s.Require().NoError(err)

	db, cleanup, err := database.NewPostgresDB(cfg)
	s.Require().NoError(err)

	s.db = db
	s.cleanup = cleanup
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func (s *UserRepositoryTestSuite) SetupTest() {
	// Clean up table before each test
	err := s.db.Exec("TRUNCATE TABLE users RESTART IDENTITY").Error
	s.Require().NoError(err)
}

func (s *UserRepositoryTestSuite) TestCreateAndGetByID() {
	repo := repository.NewUserRepository(s.db)
	ctx := context.Background()

	user := &entity.User{
		Username: "integration_test",
		Email:    "test@example.com",
	}

	// Create
	err := repo.Create(ctx, user)
	s.NoError(err)
	s.NotZero(user.ID)

	// GetByID
	found, err := repo.GetByID(ctx, user.ID)
	s.NoError(err)
	s.Equal(user.Username, found.Username)
	s.Equal(user.Email, found.Email)
}

func (s *UserRepositoryTestSuite) TestList() {
	repo := repository.NewUserRepository(s.db)
	ctx := context.Background()

	users := []entity.User{
		{Username: "user1", Email: "user1@example.com"},
		{Username: "user2", Email: "user2@example.com"},
	}

	for i := range users {
		err := repo.Create(ctx, &users[i])
		s.NoError(err)
	}

	found, err := repo.List(ctx)
	s.NoError(err)
	s.Len(found, 2)
}

func (s *UserRepositoryTestSuite) TestUpdate() {
	repo := repository.NewUserRepository(s.db)
	ctx := context.Background()

	user := &entity.User{Username: "oldname", Email: "old@example.com"}
	err := repo.Create(ctx, user)
	s.NoError(err)

	user.Username = "newname"
	err = repo.Update(ctx, user)
	s.NoError(err)

	found, err := repo.GetByID(ctx, user.ID)
	s.NoError(err)
	s.Equal("newname", found.Username)
}

func (s *UserRepositoryTestSuite) TestDelete() {
	repo := repository.NewUserRepository(s.db)
	ctx := context.Background()

	user := &entity.User{Username: "todelete", Email: "delete@example.com"}
	err := repo.Create(ctx, user)
	s.NoError(err)

	err = repo.Delete(ctx, user.ID)
	s.NoError(err)

	found, err := repo.GetByID(ctx, user.ID)
	s.Error(err)
	s.Nil(found)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
