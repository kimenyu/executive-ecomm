package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (s *Store) GetOrderWithItemsByID(orderID uuid.UUID) (*types.OrderWithItems, error) {
	query := `
		SELECT 
			o.id, o.user_id, o.total, o.status, o.address_id, o.created_at,
			oi.id, oi.product_id, oi.quantity, oi.price,
			p.name
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		LEFT JOIN products p ON oi.product_id = p.id
		WHERE o.id = $1
	`

	rows, err := s.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order types.Order
	items := []types.OrderItemDetailed{}
	firstRow := true

	for rows.Next() {
		var (
			itemID      sql.NullString
			productID   sql.NullString
			quantity    sql.NullInt32
			price       sql.NullFloat64
			productName sql.NullString
		)

		if firstRow {
			if err := rows.Scan(
				&order.ID, &order.UserID, &order.Total, &order.Status, &order.AddressID, &order.CreatedAt,
				&itemID, &productID, &quantity, &price,
				&productName,
			); err != nil {
				return nil, err
			}
			firstRow = false
		} else {
			var dummyOrderID, dummyUserID, dummyAddressID string
			var dummyTotal float64
			var dummyStatus string
			var dummyCreatedAt time.Time

			if err := rows.Scan(
				&dummyOrderID, &dummyUserID, &dummyTotal, &dummyStatus, &dummyAddressID, &dummyCreatedAt,
				&itemID, &productID, &quantity, &price,
				&productName,
			); err != nil {
				return nil, err
			}
		}

		if itemID.Valid && productID.Valid && quantity.Valid && price.Valid {
			oiID, err := uuid.Parse(itemID.String)
			if err != nil {
				return nil, fmt.Errorf("invalid order item id: %v", err)
			}
			pID, err := uuid.Parse(productID.String)
			if err != nil {
				return nil, fmt.Errorf("invalid product id: %v", err)
			}

			item := types.OrderItemDetailed{
				ID:          oiID,
				ProductID:   pID,
				ProductName: productName.String,
				Quantity:    int(quantity.Int32),
				Price:       price.Float64,
			}
			items = append(items, item)
		}
	}

	if firstRow {
		return nil, sql.ErrNoRows
	}

	return &types.OrderWithItems{
		Order: order,
		Items: items,
	}, nil
}

func (s *Store) UpdateOrder(o *types.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE orders
        SET status = $1, updated_at = $2
        WHERE id = $3
    `

	_, err := s.db.ExecContext(ctx, query, o.Status, time.Now(), o.ID)
	return err
}
