package product

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/helpers"
	"github.com/kimenyu/executive/types"
)

type Store struct {
	db *sql.DB
}

// constructor
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// create a product
func (s *Store) CreateProduct(product *types.Product) error {
	_, err := s.db.Exec(`INSERT INTO products(id, name, decsription, price, image, category_id, quantity, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, product.ID, product.Name, product.Description, product.Price, product.Image, product.CategoryID, product.Quantity, product.CreatedAt, product.UpdatedAt)

	return err
}

// get all products
func (s *Store) GetAllProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]*types.Product, 0)

	for rows.Next() {
		p, err := helpers.ScanRowsIntoProducts(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// get product by id
func (s *Store) GetProductByID(id uuid.UUID) (*types.Product, error) {
	row := s.db.QueryRow("SELECT * FROM products WHERE id=$1", id)
	return helpers.ScanRowIntoProduct(row)
}

// update a product
func (s *Store) UpdateProduct(product *types.Product) error {
	_, err := s.db.Exec(`
		UPDATE products 
		SET name = $1, 
		    description = $2, 
		    price = $3, 
		    image = $4, 
		    category_id = $5, 
		    quantity = $6, 
		    updated_at = $7
		WHERE id = $8
	`, product.Name, product.Description, product.Price, product.Image,
		product.CategoryID, product.Quantity, product.UpdatedAt, product.ID)

	return err
}

// delete product
func (s *Store) DeleteProduct(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id=$1", id)
	return err
}
