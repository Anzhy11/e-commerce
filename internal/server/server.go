package server

import (
	"net/http"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/events"
	"github.com/anzhy11/go-e-commerce/internal/interfaces"
	"github.com/anzhy11/go-e-commerce/internal/server/middlewares"
	authRoutes "github.com/anzhy11/go-e-commerce/internal/server/routes/auth"
	cartRoutes "github.com/anzhy11/go-e-commerce/internal/server/routes/cart"

	orderRoutes "github.com/anzhy11/go-e-commerce/internal/server/routes/orders"
	productRoutes "github.com/anzhy11/go-e-commerce/internal/server/routes/products"
	userRoutes "github.com/anzhy11/go-e-commerce/internal/server/routes/users"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	cfg      *config.Config
	db       *gorm.DB
	log      *zerolog.Logger
	mdw      *middlewares.Middlewares
	eventPub events.PublisherInterface
	up       interfaces.Upload
}

func New(cfg *config.Config, db *gorm.DB, log *zerolog.Logger, eventPub events.PublisherInterface, up interfaces.Upload) *Server {
	return &Server{
		cfg:      cfg,
		db:       db,
		log:      log,
		mdw:      middlewares.New(cfg),
		eventPub: eventPub,
		up:       up,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.mdw.CorsMiddleware())

	router.GET("/health", healthCheckHandler)
	router.Static("/uploads", "./uploads")

	apiGroup := router.Group("/api/v1")

	authRoutes.Setup(apiGroup, s.db, s.cfg, s.log, s.eventPub)
	userRoutes.Setup(apiGroup, s.mdw, s.db)
	productRoutes.Setup(apiGroup, s.mdw, s.db, s.cfg, s.up)
	cartRoutes.Setup(apiGroup, s.mdw, s.db)

	orderService := orderRoutes.New(apiGroup, s.mdw, s.db)
	orderService.SetupRoutes()

	return router
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
