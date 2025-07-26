package category

import (
	"database/sql"
	"github.com/kimenyu/executive/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Create a new category
func (s *Store) CreateCategory(category *types.Category) error {
	_, err := s.db.Exec(
		`INSERT INTO categories (id, name, created_at) VALUES ($1, $2, $3)`,
		category.ID, category.Name, category.CreatedAt,
	)
	return err
}

// Get all categories
func (s *Store) GetCategories() ([]*types.Category, error) {
	rows, err := s.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// Get a single category by ID
func (s *Store) GetCategoryById(id int) (*types.Category, error) {
	row := s.db.QueryRow("SELECT * FROM categories WHERE id = $1", id)
	return scanRowIntoCategory(row)
}

// Single row (QueryRow)
func scanRowIntoCategory(row *sql.Row) (*types.Category, error) {
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

// Multiple rows (Query)
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
