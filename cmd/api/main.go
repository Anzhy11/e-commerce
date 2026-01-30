package main

import (
	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/database"
	"github.com/anzhy11/go-e-commerce/internal/logger"
	"github.com/gin-gonic/gin"
)

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
		if err := mainDb.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()
	gin.SetMode(cfg.Server.GinMode)

	log.Info().Msg("Connected to database")
}
