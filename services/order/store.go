package order

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order *types.Order) error {
	_, err := s.db.Exec(`INSERT INTO orders (id, user_id, total, status, address_id, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		order.ID, order.UserID, order.Total, order.Status, order.AddressID, order.CreatedAt)
	return err
}

func (s *Store) AddOrderItem(item *types.OrderItem) error {
	_, err := s.db.Exec(`INSERT INTO order_items (id, order_id, product_id, quantity, price) 
		VALUES ($1, $2, $3, $4, $5)`,
		item.ID, item.OrderID, item.ProductID, item.Quantity, item.Price)
	return err
}

func (s *Store) GetOrdersByUser(userID uuid.UUID) ([]types.Order, error) {
	rows, err := s.db.Query(`SELECT id, user_id, total, status, address_id, created_at FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var o types.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Total, &o.Status, &o.AddressID, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
