package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mamochiro/go-microservice-template/pkg/metrics"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Use the responseWriter wrapper we already have in logger.go
		// but since it's in the same package we can just use it
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.status)

		metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
