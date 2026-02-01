package productRoutes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/interfaces"
	productHandler "github.com/anzhy11/go-e-commerce/internal/server/handlers/products"
	"github.com/anzhy11/go-e-commerce/internal/server/middlewares"
)

type productRoutes struct {
	pd productHandler.ProductHandlerInterface
}

func Setup(routeGroup *gin.RouterGroup, mdw *middlewares.Middlewares, db *gorm.DB, cfg *config.Config, up interfaces.Upload) {
	pd := productHandler.New(db, cfg, up)

	pr := &productRoutes{
		pd: pd,
	}

	prg := routeGroup.Group("/products")

	// Public routes
	prg.GET("/", pr.pd.GetProducts)
	prg.GET("/categories", pr.pd.GetCategories)
	prg.GET("/:id", pr.pd.GetProductById)

	// Protected routes
	prg.Use(mdw.Authorization())
	prg.Use(mdw.AdminAuthorization())

	prg.POST("/", pr.pd.CreateProduct)
	prg.PUT("/:id", pr.pd.UpdateProduct)
	prg.DELETE("/:id", pr.pd.DeleteProduct)
	prg.POST("/:id/upload", pr.pd.UploadProductImage)

	prg.POST("/categories", pr.pd.CreateCategory)
	prg.PUT("/categories/:id", pr.pd.UpdateCategory)
	prg.DELETE("/categories/:id", pr.pd.DeleteCategory)
}
