package review

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
)

type Store struct {
	db *sql.DB
}

// constructor
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// create a review
func (s *Store) CreateReview(review *types.Review) error {
	_, err := s.db.Exec(`INSERT INTO reviews(id, product_id, user_id, rating, comment, created_at)
				VALUES($1, $2, $3, $4, $5, $6)`, review.ID, review.ProductID, review.UserID, review.Rating, review.Comment, review.CreatedAt)

	return err
}

// delete a review by ID and user ID (ownership check)
func (s *Store) DeleteReview(id, userID uuid.UUID) error {
	res, err := s.db.Exec(`DELETE FROM reviews WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("unauthorized or review not found")
	}
	return nil
}

func (s *Store) GetReviewByID(id uuid.UUID) (*types.Review, error) {
	row := s.db.QueryRow(`SELECT id, product_id, user_id, rating, comment, created_at FROM reviews WHERE id = $1`, id)

	review := new(types.Review)
	err := row.Scan(&review.ID, &review.ProductID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt)
	return review, err
}

func (s *Store) GetReviewsByProduct(productID uuid.UUID) ([]*types.Review, error) {
	rows, err := s.db.Query(`SELECT id, product_id, user_id, rating, comment, created_at FROM reviews WHERE product_id = $1`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*types.Review
	for rows.Next() {
		r := new(types.Review)
		err := rows.Scan(&r.ID, &r.ProductID, &r.UserID, &r.Rating, &r.Comment, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}

func (s *Store) UpdateReview(review *types.Review) error {
	_, err := s.db.Exec(`UPDATE reviews SET rating = $1, comment = $2 WHERE id = $3 AND user_id = $4`, review.Rating, review.Comment, review.ID, review.UserID)
	return err
}
