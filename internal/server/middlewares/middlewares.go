package middlewares

import (
	"net/http"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Middlewares {
	return &Middlewares{
		cfg: cfg,
	}
}

func (m *Middlewares) CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
