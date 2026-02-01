package productService

import (
	"errors"

	"github.com/anzhy11/go-e-commerce/internal/dto"
	"github.com/anzhy11/go-e-commerce/internal/models"
	"github.com/anzhy11/go-e-commerce/internal/repository"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"gorm.io/gorm"
)

type ProductServiceInterface interface {
	CreateCategory(data *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories() ([]dto.CategoryResponse, error)
	UpdateCategory(categoryID uint, data *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(categoryID uint) error
	CreateProduct(data *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginatedMeta, error)
	GetProductById(productID uint) (*dto.ProductResponse, error)
	UpdateProduct(productID uint, data *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(productID uint) error
	AddProductImage(productID uint, url, altText string) error
}

type productService struct {
	db          *gorm.DB
	productRepo *repository.ProductRepository
}

func New(db *gorm.DB) ProductServiceInterface {
	return &productService{
		db:          db,
		productRepo: repository.NewProductRepo(db),
	}
}

// Category
func (s *productService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.productRepo.CreateCategory(&category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}

func (s *productService) GetCategories() ([]dto.CategoryResponse, error) {
	category, err := s.productRepo.GetAllCategories()
	if err != nil {
		return nil, err
	}

	categoryResponses := make([]dto.CategoryResponse, len(category))
	for i := range category {
		categoryResponses[i] = dto.CategoryResponse{
			ID:          category[i].ID,
			Name:        category[i].Name,
			Description: category[i].Description,
			IsActive:    category[i].IsActive,
		}
	}
	return categoryResponses, nil
}

func (s *productService) UpdateCategory(categoryID uint, data *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.productRepo.GetCategoryById(categoryID)
	if err != nil {
		return nil, err
	}

	category.Name = data.Name
	category.Description = data.Description
	if data.IsActive {
		category.IsActive = data.IsActive
	}

	if err := s.productRepo.UpdateCategory(category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}

func (s *productService) DeleteCategory(categoryID uint) error {
	return s.productRepo.DeleteCategory(categoryID)
}

// Product
func (s *productService) CreateProduct(data *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Stock:       data.Stock,
		SKU:         data.SKU,
		CategoryID:  data.CategoryID,
	}

	if err := s.productRepo.CreateProduct(&product); err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
	}, nil
}

func (s *productService) GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginatedMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, total, err := s.productRepo.GetProducts(offset, limit)
	if err != nil {
		return nil, nil, err
	}

	productResponses := make([]dto.ProductResponse, len(products))
	for i := range products {
		productResponses[i] = *s.generateProductResponse(&products[i])
	}

	totalPages := int((total + int64(limit)) / int64(limit))
	meta := &utils.PaginatedMeta{
		Page:       page,
		PageSize:   limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return productResponses, meta, nil
}

func (s *productService) GetProductById(productID uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductById(productID)
	if err != nil {
		return nil, err
	}

	return s.generateProductResponse(product), nil
}

func (s *productService) UpdateProduct(productID uint, data *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductById(productID)
	if err != nil {
		return nil, err
	}

	product.Name = data.Name
	product.Description = data.Description
	product.Price = data.Price
	product.Stock = data.Stock
	product.SKU = data.SKU
	product.CategoryID = data.CategoryID
	if data.IsActive {
		product.IsActive = data.IsActive
	}

	if err := s.productRepo.UpdateProduct(product); err != nil {
		return nil, err
	}

	return s.generateProductResponse(product), nil
}

func (s *productService) DeleteProduct(productID uint) error {
	return s.productRepo.DeleteProduct(productID)
}

func (s *productService) AddProductImage(productID uint, url, altText string) error {
	count := s.productRepo.CountProductImage(productID)
	if count >= 10 {
		return errors.New("maximum number of images reached")
	}

	image := models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   altText,
	}

	if count == 0 {
		image.IsPrimary = true
	}

	if err := s.productRepo.UploadProductImage(&image); err != nil {
		return err
	}

	return nil
}

// Helper
func (s *productService) generateProductResponse(product *models.Product) *dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
		}
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			IsActive:    product.Category.IsActive,
		},
		Images: images,
	}
}
