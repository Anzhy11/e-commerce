package repository

import (
	"github.com/anzhy11/go-e-commerce/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Product category
func (r *ProductRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *ProductRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductRepository) GetCategoryById(categoryID uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ProductRepository) UpdateCategory(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *ProductRepository) DeleteCategory(categoryID uint) error {
	return r.db.Delete(categoryID).Error
}

// Product

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetProducts(offset, limit int) ([]models.Product, int64, error) {
	var products []models.Product

	total := r.GetProductsCount()

	if err := r.db.Preload("Category").Preload("Images").
		Where("is_active = ?", true).
		Offset(offset).Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepository) GetProductsCount() int64 {
	var total int64
	r.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total)

	return total
}

func (r *ProductRepository) GetProductById(productID uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").First(&product, productID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) DeleteProduct(productID uint) error {
	return r.db.Delete(productID).Error
}

func (r *ProductRepository) UploadProductImage(image *models.ProductImage) error {
	return r.db.Create(&image).Error
}

func (r *ProductRepository) DeleteProductImage(imageID uint) error {
	return r.db.Delete(&models.ProductImage{}, imageID).Error
}

func (r *ProductRepository) CountProductImage(productID uint) int64 {
	var count int64
	r.db.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count)
	return count
}
