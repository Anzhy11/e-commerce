package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	IsActive    bool   `json:"is_active" binding:"omitempty"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price" binding:"omitempty"`
	Stock       int     `json:"stock" binding:"omitempty"`
	SKU         string  `json:"sku" binding:"omitempty"`
	IsActive    bool    `json:"is_active" binding:"omitempty"`
}

type ProductResponse struct {
	ID          uint             `json:"id"`
	CategoryID  uint             `json:"category_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Stock       int              `json:"stock"`
	SKU         string           `json:"sku"`
	IsActive    bool             `json:"is_active"`
	Category    CategoryResponse `json:"category"`
	Images      []ImageResponse  `json:"images"`
	UpdatedAt   string           `json:"updated_at"`
}

type ImageResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
}
