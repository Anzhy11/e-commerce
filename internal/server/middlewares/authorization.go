package middlewares

import (
	"slices"
	"strings"

	"github.com/anzhy11/go-e-commerce/internal/models"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Middlewares) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		payload, err := utils.VerifyToken(tokenString, s.cfg.JWT.Secret)
		if err != nil {
			utils.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("user_id", payload.UserID)
		c.Set("email", payload.Email)
		c.Set("role", payload.Role)
		c.Next()
	}
}

func (s *Middlewares) AdminAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			utils.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		if role != string(models.RoleAdmin) {
			utils.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *Middlewares) RoleAuthorization(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			utils.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		if slices.Contains(roles, role) {
			c.Next()
			return
		}

		utils.Forbidden(c, "Forbidden")
		c.Abort()
	}
}
