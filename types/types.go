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
	CreateUser(user *User) error
}

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
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

type CreateCategoryPayload struct {
	Name string `json:"name" validate:"required"`
}
type CategoryStore interface {
	CreateCategory(category *Category) error
	GetCategories() ([]*Category, error)
	GetCategoryById(id uuid.UUID) (*Category, error)
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Image       string    `json:"image"`
	CategoryID  uuid.UUID `json:"category_id"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// used in the http layer only(to handler user input)
type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required"`
	Image       string  `json:"image"`
	CategoryID  string  `json:"category_id"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type ProductStore interface {
	CreateProduct(product *Product) error
	GetProductByID(id uuid.UUID) (*Product, error)
	GetAllProducts() ([]*Product, error)
	DeleteProduct(id uuid.UUID) error
	UpdateProduct(product *Product) error
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
