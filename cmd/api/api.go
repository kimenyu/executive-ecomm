package api

import (
	"database/sql"
	"github.com/kimenyu/executive/services/user"
	"log"
	"net/http"

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

	// Middleware like CORS, logging, etc. can go here (optional)

	router.Route("/api/v1", func(r chi.Router) {
		// Setup stores
		userStore := user.NewStore(s.db)
		productStore := product.NewStore(s.db)

		// Setup handlers
		userHandler := user.NewHandler(userStore)
		productHandler := product.NewHandler(productStore)

		// Register public routes
		userHandler.RegisterRoutes(r)
		productHandler.RegisterRoutes(r)

	})

	log.Printf("Server listening on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
