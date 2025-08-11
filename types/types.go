package types

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
	GetUserByID(id uuid.UUID) (*User, error)
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

type CartStore interface {
	GetCartByUserID(userID uuid.UUID) (*Cart, error)
	CreateCart(cart *Cart) error
	AddCartItem(item *CartItem) error
	GetCartItems(cartID uuid.UUID) ([]CartItem, error)
}

type AddToCartPayload struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
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
	Price     float64   `json:"price"`
}

type CreateOrderPayload struct {
	Items []CreateOrderItemDTO `json:"items" validate:"required,dive"`
	Total float64              `json:"total" validate:"required,gt=0"`
}

type CreateOrderItemDTO struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	Price     float64   `json:"price" validate:"required,gt=0"`
}

type OrderItemDetailed struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
}

type OrderWithItems struct {
	Order Order               `json:"order"`
	Items []OrderItemDetailed `json:"items"`
}

type UpdateOrderPayload struct {
	Status string `json:"status" validate:"required,oneof=pending paid shipped completed cancelled"`
}

type OrderStore interface {
	CreateOrder(order *Order) error
	AddOrderItem(item *OrderItem) error
	GetOrdersByUser(userID uuid.UUID) ([]Order, error)
	GetOrderWithItemsByID(orderID uuid.UUID) (*OrderWithItems, error)
	UpdateOrder(order *Order) error
}
type Payment struct {
	ID                uuid.UUID       `json:"id"`
	OrderID           uuid.UUID       `json:"order_id"`
	Amount            float64         `json:"amount"`
	Provider          string          `json:"provider"`
	Status            string          `json:"status"`
	CheckoutRequestID string          `json:"checkout_request_id"`
	MerchantRequestID string          `json:"merchant_request_id"`
	MpesaReceipt      string          `json:"mpesa_receipt"`
	Phone             string          `json:"phone"`
	Metadata          json.RawMessage `json:"metadata"`
	CreatedAt         time.Time       `json:"created_at"`
}

type Review struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int       `json:"rating"` // 1 to 5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateReviewPayload struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type ReviewStore interface {
	CreateReview(review *Review) error
	GetReviewByID(id uuid.UUID) (*Review, error)
	GetReviewsByProduct(productID uuid.UUID) ([]*Review, error)
	UpdateReview(review *Review) error
	DeleteReview(id uuid.UUID, userID uuid.UUID) error
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

type AddressStore interface {
	CreateAddress(address *Address) error
	GetAddress(userID uuid.UUID) (*Address, error)
	UpdateAddress(address *Address) error
}

type CreateAddressPayload struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	City    string `json:"city" validate:"required"`
	Country string `json:"country" validate:"required"`
	ZipCode string `json:"zip_code" validate:"required"`
}
