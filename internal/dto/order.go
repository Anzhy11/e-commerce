package dto

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	Total     int                `json:"total"`
	CartItems []CartItemResponse `json:"cart_items"`
}

type CartItemResponse struct {
	ID       uint            `json:"id"`
	Quantity int             `json:"quantity"`
	Subtotal int             `json:"subtotal"`
	Product  ProductResponse `json:"product"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount int                 `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   string              `json:"created_at"`
}

type OrderItemResponse struct {
	ID       uint            `json:"id"`
	Quantity int             `json:"quantity"`
	Price    int             `json:"price"`
	Product  ProductResponse `json:"product"`
}
