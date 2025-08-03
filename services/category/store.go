package category

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/helpers"
	"github.com/kimenyu/executive/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// create a category
func (s *Store) CreateCategory(category *types.Category) error {
	_, err := s.db.Exec(`INSERT INTO categories(id, name, created_at) VALUES($1, $2, $3)`, category.ID, category.Name, category.CreatedAt)
	return err
}

// get all categories
func (s *Store) GetCategories() ([]*types.Category, error) {
	rows, err := s.db.Query("SELECT  * from categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// slice to hold categories
	categories := make([]*types.Category, 0)

	for rows.Next() {
		c, err := helpers.ScanRowsIntoCategories(rows)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

// get category by id
func (s *Store) GetCategoryById(id uuid.UUID) (*types.Category, error) {
	row := s.db.QueryRow("SELECT * FROM categories WHERE id=$1", id)
	return helpers.ScanRowIntoCategory(row)
}
