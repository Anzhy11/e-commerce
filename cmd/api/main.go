package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/database"
	"github.com/anzhy11/go-e-commerce/internal/events"
	"github.com/anzhy11/go-e-commerce/internal/interfaces"
	"github.com/anzhy11/go-e-commerce/internal/logger"
	"github.com/anzhy11/go-e-commerce/internal/providers"
	"github.com/anzhy11/go-e-commerce/internal/server"
	"github.com/gin-gonic/gin"
)

// @title E-Commerce API
// @version 1.0
// @description A modern e-commerce API built with Go, Gin, and GORM
// @termsOfService http://swagger.io/terms/

// @contact.name   Ahmed Adebayo
// @contact.url    http://linkedin.com/in/anzhy11
// @contact.email  codemastery.web@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemas http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	mainDb, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database connection")
	}

	defer func() {
		if dbErr := mainDb.Close(); dbErr != nil {
			log.Error().Err(dbErr).Msg("Failed to close database connection")
		}
	}()
	gin.SetMode(cfg.Server.GinMode)

	var up interfaces.Upload
	if cfg.Upload.Provider == "s3" {
		up = providers.NewS3UploadProvider(cfg)
	} else {
		up = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}

	ctx := context.Background()
	eventPub, err := events.NewEventPublisher(ctx, &cfg.AWS)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create event publisher")
	}

	srv := server.New(cfg, db, log, eventPub, up)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Starting http server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited properly")
}
