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

// scan single category
func ScanRowIntoCategory(row *sql.Row) (*types.Category, error) {
	category := new(types.Category)

	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return category, nil

}

// scan mulitplecategory
func ScanRowsIntoCategories(rows *sql.Rows) (*types.Category, error) {
	categories := new(types.Category)

	err := rows.Scan(
		&categories.ID,
		&categories.Name,
		&categories.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return categories, nil

}

func ScanRowIntoAddress(row *sql.Row) (*types.Address, error) {

	address := new(types.Address)

	err := row.Scan(
		&address.ID,
		&address.UserID,
		&address.Line1,
		&address.Line2,
		&address.City,
		&address.Country,
		&address.ZipCode,
		&address.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return address, nil
}
