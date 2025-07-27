package product

import (
	"database/sql"
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
	_, err := s.db.Exec(`
		INSERT INTO products (id, name, description, price, image, category_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		product.ID, product.Name, product.Description, product.Price,
		product.Image, product.CategoryID, product.Quantity,
		product.CreatedAt, product.UpdatedAt,
	)
	return err
}

// get a product by ID
func (s *Store) GetProductByID(productID int) (*types.Product, error) {
	row := s.db.QueryRow("SELECT * FROM products WHERE id = $1", productID)
	return scanRowIntoProduct(row)
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
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// delete a product by ID
func (s *Store) DeleteProduct(id int) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}

// update a product by ID
func (s *Store) UpdateProduct(product *types.Product) error {
	_, err := s.db.Exec(`
		UPDATE products 
		SET name = $1, description = $2, price = $3, image = $4, category_id = $5, quantity = $6, updated_at = $7 
		WHERE id = $8`,
		product.Name, product.Description, product.Price,
		product.Image, product.CategoryID, product.Quantity,
		product.UpdatedAt, product.ID,
	)
	return err
}

// single row scan
func scanRowIntoProduct(row *sql.Row) (*types.Product, error) {
	product := new(types.Product)
	err := row.Scan(
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

// multiple rows scan
func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
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
