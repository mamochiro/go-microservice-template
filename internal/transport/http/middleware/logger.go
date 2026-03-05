package middleware

import (
	"net/http"
	"time"

	"github.com/mamochiro/go-microservice-template/pkg/logger"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		latency := time.Since(start)

		logger.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.status),
			zap.Duration("latency", latency),
			zap.String("ip", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.String("correlation_id", r.Header.Get("X-Correlation-ID")),
		)
	})
}
