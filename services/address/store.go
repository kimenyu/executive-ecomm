package address

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

func (s *Store) CreateAddress(address *types.Address) error {
	_, err := s.db.Exec(`INSERT INTO addresses(id, user_id, line1, line2, city, country, zip_code, created_at) VALUES($1, $2, $3, $4, $5, $6,  $7, $8)
			`, address.ID, address.UserID, address.Line1, address.Line2, address.City, address.Country, address.ZipCode, address.CreatedAt)

	return err
}

func (s *Store) GetAddress(userID uuid.UUID) (*types.Address, error) {
	row := s.db.QueryRow(`SELECT id, user_id, line1, line2, city, country, zip_code, created_at FROM addresses WHERE user_id=$1`, userID)
	return helpers.ScanRowIntoAddress(row)
}

func (s *Store) UpdateAddress(address *types.Address) error {
	_, err := s.db.Exec(`
		UPDATE addresses
		SET line1 = $1,
		    line2 = $2,
		    city = $3,
		    country = $4,
		    zip_code = $5
		WHERE id = $8`,
		address.Line1, address.Line2, address.City, address.Country, address.ZipCode, address.ID)
	return err
}
