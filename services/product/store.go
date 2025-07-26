package product

import (
	"database/sql"
	"github.com/kimenyu/executive/types"
)

type Store struct {
	db *sql.DB
}

// constructor for the above store
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// data access to create a product
func (s *Store) CreateProduct(product *types.Product) error {
	_, err := s.db.Exec(`INSERT INTO products(id, name, description, price, image, category_id, quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		product.ID, product.Name, product.Description, product.Price, product.Image, product.CategoryID, product.Quantity, product.CreatedAt, product.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// get a product by id
func (s *Store) GetProductByID(productID int) (*types.Product, error) {
	row := s.db.QueryRow("SELECT * FROM products WHERE id = $1", productID)
	return scanRowsIntoProduct(row)
}

func scanRowsIntoProduct(rows *sql.Row) (*types.Product, error) {

	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Image,
		&product.CategoryID,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}
