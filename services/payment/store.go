package payment

import (
	"database/sql"

	"github.com/kimenyu/executive/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

func (s *Store) CreatePayment(p *types.Payment) error {
	_, err := s.db.Exec(
		`INSERT INTO payments (id, order_id, amount, provider, status, checkout_request_id, merchant_request_id, mpesa_receipt, phone, metadata, created_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		p.ID, p.OrderID, p.Amount, p.Provider, p.Status, p.CheckoutRequestID, p.MerchantRequestID, p.MpesaReceipt, p.Phone, p.Metadata, p.CreatedAt,
	)
	return err
}

func (s *Store) GetPaymentByCheckoutID(checkout string) (*types.Payment, error) {
	row := s.db.QueryRow(`SELECT id, order_id, amount, provider, status, checkout_request_id, merchant_request_id, mpesa_receipt, phone, metadata, created_at FROM payments WHERE checkout_request_id=$1`, checkout)
	var p types.Payment
	var raw []byte
	if err := row.Scan(&p.ID, &p.OrderID, &p.Amount, &p.Provider, &p.Status, &p.CheckoutRequestID, &p.MerchantRequestID, &p.MpesaReceipt, &p.Phone, &raw, &p.CreatedAt); err != nil {
		return nil, err
	}
	p.Metadata = raw
	return &p, nil
}
