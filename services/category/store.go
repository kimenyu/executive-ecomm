package category

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateCategory(category *types.Category) error {
	_, err := s.db.Exec("INSERT INTO categories (name) VALUES ($1), category.name")

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetCategories() ([]*types.Category, error) {
	rows, err := s.db.Query("SELECT * FROM categories")

	if err != nil {
		return nil, err
	}

	// create an empty category slice
	categories := make([]*types.Category, 0)

	for rows.Next() {
		c, err := scanRowsIntoCategory(rows)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func (s *Store) GetCategoryById(id int) (*types.Category, error) {
	rows, err := s.db.Query("SELECT * FROM categories WHERE id=$1", id)

	if err != nil {
		return nil, err
	}

	// create an empty category struct
	c := new(types.Category)
	for rows.Next() {
		c, err := scanRowsIntoCategory(rows)

		if err != nil {
			return nil, err
		}
	}

	if c.ID != uuid.Nil {
		return nil, fmt.Errorf("Category not found")
	}

	return c, nil
}

// helper function to convert the db row to category struct
func scanRowsIntoCategory(rows *sql.Rows) (*types.Category, error) {

	category := new(types.Category)

	err := rows.Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}
