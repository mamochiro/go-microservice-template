package database

import (
	"fmt"
	"log"

	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, func(), error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.Port,
		cfg.Postgres.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	// Auto Migration for example
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		sqlDB.Close()
		log.Println("PostgreSQL connection closed")
	}

	log.Println("Connected to PostgreSQL database and migrated")
	return db, cleanup, nil
}
