package productHandler

import (
	"strconv"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/dto"
	"github.com/anzhy11/go-e-commerce/internal/interfaces"
	productService "github.com/anzhy11/go-e-commerce/internal/services/products"
	uploadService "github.com/anzhy11/go-e-commerce/internal/services/upload"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandlerInterface interface {
	CreateCategory(c *gin.Context)
	GetCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	CreateProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	GetProductById(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	UploadProductImage(c *gin.Context)
}

type productHandler struct {
	pd  productService.ProductServiceInterface
	db  *gorm.DB
	us  *uploadService.UploadService
	cfg *config.Config
}

func New(db *gorm.DB, cfg *config.Config, up interfaces.Upload) ProductHandlerInterface {
	us := uploadService.NewUploadService(up)

	return &productHandler{
		pd:  productService.New(db),
		db:  db,
		cfg: cfg,
		us:  us,
	}
}

// Category
func (h *productHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	category, err := h.pd.CreateCategory(&req)
	if err != nil {
		utils.InternalServerError(c, "failed to create category", err)
		return
	}

	utils.SuccessResponse(c, "Category created successfully", category)
}

func (h *productHandler) GetCategories(c *gin.Context) {
	categories, err := h.pd.GetCategories()
	if err != nil {
		utils.InternalServerError(c, "failed to get categories", err)
		return
	}

	utils.SuccessResponse(c, "Categories fetched successfully", categories)
}

func (h *productHandler) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "invalid category id", err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	category, catErr := h.pd.UpdateCategory(uint(categoryID), &req)
	if catErr != nil {
		utils.InternalServerError(c, "failed to update category", catErr)
		return
	}

	utils.SuccessResponse(c, "Category updated successfully", category)
}

func (h *productHandler) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "invalid category id", err)
		return
	}

	if err := h.pd.DeleteCategory(uint(categoryID)); err != nil {
		utils.InternalServerError(c, "failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "Category deleted successfully", nil)
}

// Product
func (h *productHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	product, err := h.pd.CreateProduct(&req)
	if err != nil {
		utils.InternalServerError(c, "failed to create product", err)
		return
	}

	utils.SuccessResponse(c, "Product created successfully", product)
}

func (h *productHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, meta, err := h.pd.GetProducts(page, limit)
	if err != nil {
		utils.InternalServerError(c, "failed to get products", err)
		return
	}

	utils.SuccessResponse(c, "Products fetched successfully", gin.H{
		"products": products,
		"meta":     *meta,
	})
}

func (h *productHandler) GetProductById(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "invalid product id", err)
		return
	}

	product, err := h.pd.GetProductById(uint(productID))
	if err != nil {
		utils.InternalServerError(c, "failed to get product", err)
		return
	}

	utils.SuccessResponse(c, "Product fetched successfully", product)
}

func (h *productHandler) UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "invalid product id", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	product, prdErr := h.pd.UpdateProduct(uint(productID), &req)
	if prdErr != nil {
		utils.InternalServerError(c, "failed to update product", prdErr)
		return
	}

	utils.SuccessResponse(c, "Product updated successfully", product)
}

func (h *productHandler) DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "invalid product id", err)
		return
	}

	if err := h.pd.DeleteProduct(uint(productID)); err != nil {
		utils.InternalServerError(c, "failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product deleted successfully", nil)
}

func (h *productHandler) UploadProductImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "invalid product id", err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequest(c, "invalid image", err)
		return
	}

	url, err := h.us.UploadProductImage(uint(productID), file)
	if err != nil {
		utils.InternalServerError(c, "failed to upload image", err)
		return
	}

	if err := h.pd.AddProductImage(uint(productID), url, ""); err != nil {
		utils.InternalServerError(c, "failed to upload image", err)
		return
	}

	utils.SuccessResponse(c, "Image uploaded successfully", nil)
}
