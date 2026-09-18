package domain

import "context"

// Menu Item Model
type MenuItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"` // e.g., "Starters", "Main Course", "Beverages"
	ImageURL string  `json:"image_url"`
}

// Order Model
type Order struct {
	ID          int     `json:"id"`
	TableNumber int     `json:"table_number"`
	Items       string  `json:"items"` // Format: "2x Burger, 1x Coke"
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"` // "pending" or "completed"
}

// Repository (Database) Interface
type MenuRepository interface {
	GetMenu(ctx context.Context) ([]MenuItem, error)
	CreateOrder(ctx context.Context, order *Order) error
	GetActiveOrders(ctx context.Context) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int, status string) error
}

// Usecase (Business Logic) Interface
type MenuUsecase interface {
	FetchMenu(ctx context.Context) ([]MenuItem, error)
	PlaceOrder(ctx context.Context, order *Order) error
	FetchActiveOrders(ctx context.Context) ([]Order, error)
	MarkOrderCompleted(ctx context.Context, orderID int) error
}
