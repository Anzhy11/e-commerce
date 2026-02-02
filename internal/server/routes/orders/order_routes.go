package orderRoutes

import (
	orderHandler "github.com/anzhy11/go-e-commerce/internal/server/handlers/orders"
	"github.com/anzhy11/go-e-commerce/internal/server/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type orderRoutes struct {
	routeGroup   *gin.RouterGroup
	mdw          *middlewares.Middlewares
	orderHandler orderHandler.OrderHandlerInterface
}

func New(routeGroup *gin.RouterGroup, mdw *middlewares.Middlewares, db *gorm.DB) *orderRoutes {
	return &orderRoutes{
		routeGroup:   routeGroup,
		mdw:          mdw,
		orderHandler: orderHandler.New(db),
	}
}

func (o *orderRoutes) SetupRoutes() {
	orderGroup := o.routeGroup.Group("/orders")
	orderGroup.Use(o.mdw.Authorization())
	orderGroup.POST("/", o.orderHandler.CreateOrder)
	orderGroup.GET("/", o.orderHandler.GetOrders)
	orderGroup.GET("/:id", o.orderHandler.GetOrder)
}
