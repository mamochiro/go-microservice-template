package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	status := "UP"
	details := make(map[string]string)

	// Check Redis
	if err := h.redis.Ping(ctx).Err(); err != nil {
		status = "DOWN"
		details["redis"] = "disconnected"
	} else {
		details["redis"] = "ok"
	}

	// Check Database
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		status = "DOWN"
		details["database"] = "disconnected"
	} else {
		details["database"] = "ok"
	}

	w.Header().Set("Content-Type", "application/json")
	if status == "DOWN" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"details": details,
		"time":    time.Now().Format(time.RFC3339),
	})
}
