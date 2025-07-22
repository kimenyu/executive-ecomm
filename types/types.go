package types

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	CategoryID  uuid.UUID `json:"category_id"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Cart struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CartItem struct {
	ID        uuid.UUID `json:"id"`
	CartID    uuid.UUID `json:"cart_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"` // pending, paid, shipped, cancelled
	AddressID uuid.UUID `json:"address_id"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"` // price at the time of order
}

type Payment struct {
	ID      uuid.UUID `json:"id"`
	OrderID uuid.UUID `json:"order_id"`
	Method  string    `json:"method"` // e.g., mpesa, card
	Status  string    `json:"status"` // e.g., success, failed
	PaidAt  time.Time `json:"paid_at"`
}

type Review struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int       `json:"rating"` // 1 to 5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type Address struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Line1     string    `json:"line1"`
	Line2     string    `json:"line2"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	ZipCode   string    `json:"zip_code"`
	CreatedAt time.Time `json:"created_at"`
}
