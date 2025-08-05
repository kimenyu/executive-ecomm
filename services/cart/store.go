package cart

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

func (s *Store) GetCartByUserID(userID uuid.UUID) (*types.Cart, error) {
	row := s.db.QueryRow(`SELECT id, user_id, created_at FROM carts WHERE user_id = $1`, userID)

	var cart types.Cart
	err := row.Scan(&cart.ID, &cart.UserID, &cart.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (s *Store) CreateCart(cart *types.Cart) error {
	_, err := s.db.Exec(`INSERT INTO carts(id, user_id, created_at) VALUES($1, $2, $3)`,
		cart.ID, cart.UserID, cart.CreatedAt)
	return err
}

func (s *Store) AddCartItem(item *types.CartItem) error {
	_, err := s.db.Exec(`INSERT INTO cart_items(id, cart_id, product_id, quantity, created_at)
		VALUES($1, $2, $3, $4, $5)`,
		item.ID, item.CartID, item.ProductID, item.Quantity, item.CreatedAt)
	return err
}

func (s *Store) GetCartItems(cartID uuid.UUID) ([]types.CartItem, error) {
	rows, err := s.db.Query(`SELECT id, cart_id, product_id, quantity, created_at FROM cart_items WHERE cart_id = $1`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.CartItem
	for rows.Next() {
		var item types.CartItem
		if err := rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
