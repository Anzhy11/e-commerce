package cartRoutes

import (
	cartHandler "github.com/anzhy11/go-e-commerce/internal/server/handlers/cart"
	"github.com/anzhy11/go-e-commerce/internal/server/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type cartRoutes struct {
	cartHandler cartHandler.CartHandlerInterface
}

func Setup(routeGroup *gin.RouterGroup, mdw *middlewares.Middlewares, db *gorm.DB) {
	cr := &cartRoutes{
		cartHandler: cartHandler.NewCartHandler(db),
	}

	crg := routeGroup.Group("/cart")
	crg.Use(mdw.Authorization())
	crg.GET("/", cr.cartHandler.GetCartByUserID)
	crg.POST("/", cr.cartHandler.AddToCart)
	crg.PUT("/:id", cr.cartHandler.UpdateCartItem)
	crg.DELETE("/:id", cr.cartHandler.RemoveFromCart)
}
