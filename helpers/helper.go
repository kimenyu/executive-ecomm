package helpers

import (
	"database/sql"
	"github.com/kimenyu/executive/types"
)

// single row
func ScanRowIntoProduct(row *sql.Row) (*types.Product, error) {

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

// multiple rows
func ScanRowsIntoProducts(rows *sql.Rows) (*types.Product, error) {

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
