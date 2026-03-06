package health

import (
	"context"
	"net/http"
	"time"

	"github.com/mamochiro/go-microservice-template/internal/transport/http/handler"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHandler(db *gorm.DB, redis *redis.Client) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
	}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	status := "UP"
	details := make(map[string]string)

	if err := h.redis.Ping(ctx).Err(); err != nil {
		status = "DOWN"
		details["redis"] = "disconnected"
	} else {
		details["redis"] = "ok"
	}

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		status = "DOWN"
		details["database"] = "disconnected"
	} else {
		details["database"] = "ok"
	}

	statusCode := http.StatusOK
	if status == "DOWN" {
		statusCode = http.StatusServiceUnavailable
	}

	handler.RespondJSON(w, statusCode, map[string]interface{}{
		"status":  status,
		"details": details,
		"time":    time.Now().Format(time.RFC3339),
	})
}
