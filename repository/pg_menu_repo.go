package repository

import (
	"context"
	"database/sql"
	"restaurant-qr/domain"
)

type pgMenuRepo struct {
	DB *sql.DB
}

// Constructor
func NewPgMenuRepository(db *sql.DB) domain.MenuRepository {
	return &pgMenuRepo{DB: db}
}

// 1. Menu Items Fetch karna
func (r *pgMenuRepo) GetMenu(ctx context.Context) ([]domain.MenuItem, error) {
	query := `SELECT id, name, price, category, COALESCE(image_url, '') FROM menu_items`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var item domain.MenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Category, &item.ImageURL); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// 2. Naya Order Database me daalna
func (r *pgMenuRepo) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `INSERT INTO orders (table_number, items, total_amount, status) VALUES ($1, $2, $3, $4) RETURNING id`
	// Return hui ID ko order object me wapas set kar rahe hain
	return r.DB.QueryRowContext(ctx, query, order.TableNumber, order.Items, order.TotalAmount, order.Status).Scan(&order.ID)
}

// 3. Kitchen ke liye Pending Orders nikalna
func (r *pgMenuRepo) GetActiveOrders(ctx context.Context) ([]domain.Order, error) {
	// Puraane orders upar dikhe isliye ORDER BY id ASC
	query := `SELECT id, table_number, items, total_amount, status FROM orders WHERE status = 'pending' ORDER BY id ASC`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.TableNumber, &o.Items, &o.TotalAmount, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// 4. Order ko "completed" mark karna
func (r *pgMenuRepo) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, status, orderID)
	return err
}
