package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	redis *redis.Client
}

func NewHealthHandler(redis *redis.Client) *HealthHandler {
	return &HealthHandler{redis: redis}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := "UP"
	if err := h.redis.Ping(context.Background()).Err(); err != nil {
		status = "DOWN (redis connection lost)"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}
