// @title Executive API
// @version 1.0
// @description This is a REST API for the Executive eCommerce platform.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/kimenyu/executive/services/cart"
	"github.com/kimenyu/executive/services/category"
	"github.com/kimenyu/executive/services/order"
	"github.com/kimenyu/executive/services/review"
	"github.com/kimenyu/executive/services/user"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/kimenyu/executive/services/product"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewRouter()
	router.Get("/swagger/*", httpSwagger.WrapHandler)
	router.Route("/api/v1", func(r chi.Router) {
		// Setup stores
		userStore := user.NewStore(s.db)
		productStore := product.NewStore(s.db)
		categoryStore := category.NewStore(s.db)
		reviewStore := review.NewStore(s.db)
		cartStore := cart.NewStore(s.db)
		orderStore := order.NewStore(s.db)

		// Setup handlers
		userHandler := user.NewHandler(userStore)
		productHandler := product.NewHandler(productStore)
		categoryHandler := category.NewHandler(categoryStore)
		reviewHandler := review.NewHandler(reviewStore)
		cartHandler := cart.NewHandler(cartStore)
		orderHandler := order.NewHandler(orderStore)

		// Register public routes
		userHandler.RegisterRoutes(r)
		productHandler.RegisterRoutes(r)
		categoryHandler.RegisterRoutes(r)
		reviewHandler.RegisterRoutes(r)
		cartHandler.RegisterRoutes(r)
		orderHandler.RegisterRoutes(r)

	})

	log.Printf("Server listening on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
