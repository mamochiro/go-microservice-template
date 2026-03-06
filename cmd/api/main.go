package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mamochiro/go-microservice-template/internal/app"
	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/pkg/logger"
	"github.com/mamochiro/go-microservice-template/pkg/telemetry"
	"go.uber.org/zap"
)

// @title           Go Microservice Template API
// @version         1.0
// @description     This is a sample microservice server using clean architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and then your token.
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Logger
	logger.Init(cfg.App.Env)

	// Initialize OpenTelemetry
	tp, err := telemetry.InitTracer(cfg.App.Name)
	if err != nil {
		logger.Fatal("Failed to initialize telemetry", zap.Error(err))
	}
	defer telemetry.ShutdownTracer(tp)

	// Initialize application using Wire
	router, cleanup, err := app.InitializeApp(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Setup server
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(fmt.Sprintf("Starting server on port %s", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Listen and serve error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline to wait for active requests to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	// Call Wire cleanup (closes DB and Redis) AFTER the server stops accepting requests
	cleanup()

	logger.Info("Server exited gracefully")
}
