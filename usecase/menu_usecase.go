package usecase

import (
	"context"
	"errors"
	"restaurant-qr/domain"
)

type menuUsecase struct {
	repo domain.MenuRepository
}

// Constructor
func NewMenuUsecase(r domain.MenuRepository) domain.MenuUsecase {
	return &menuUsecase{repo: r}
}

// 1. Menu fetch karna
func (u *menuUsecase) FetchMenu(ctx context.Context) ([]domain.MenuItem, error) {
	return u.repo.GetMenu(ctx)
}

// 2. Order place karna (business logic applied)
func (u *menuUsecase) PlaceOrder(ctx context.Context, order *domain.Order) error {
	// Validation: Table number aur items empty nahi hone chahiye
	if order.TableNumber <= 0 || order.Items == "" {
		return errors.New("invalid order details")
	}

	// Naya order hamesha 'pending' status se start hoga
	order.Status = "pending"

	return u.repo.CreateOrder(ctx, order)
}

// 3. Kitchen ke liye pending orders lana
func (u *menuUsecase) FetchActiveOrders(ctx context.Context) ([]domain.Order, error) {
	return u.repo.GetActiveOrders(ctx)
}

// 4. Order complete mark karna
func (u *menuUsecase) MarkOrderCompleted(ctx context.Context, orderID int) error {
	if orderID <= 0 {
		return errors.New("invalid order ID")
	}
	return u.repo.UpdateOrderStatus(ctx, orderID, "completed")
}
