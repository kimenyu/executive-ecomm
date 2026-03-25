// cmd/api/api.go
// @title Executive API
// @version 1.0
// @description REST API for Executive eCommerce platform.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package api

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/redis/go-redis/v9"
	"github.com/go-redis/redis_rate/v10"

	"github.com/google/uuid"

	"github.com/kimenyu/executive/internal/logging"
	"github.com/kimenyu/executive/services/address"
	"github.com/kimenyu/executive/services/cart"
	"github.com/kimenyu/executive/services/category"
	"github.com/kimenyu/executive/services/order"
	"github.com/kimenyu/executive/services/payment"
	"github.com/kimenyu/executive/services/product"
	"github.com/kimenyu/executive/services/review"
	"github.com/kimenyu/executive/services/user"
	"github.com/kimenyu/executive/types"
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

// ClientKey for rate limiting (user ID if logged in, else IP)
func ClientKey(r *http.Request) string {
	if uid := types.UserIDFromContext(r.Context()); uid != uuid.Nil {
		return "uid:" + uid.String()
	}
	return "ip:" + RealIP(r)
}

// RealIP returns client IP
func RealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

func (s *APIServer) Run() error {
	logging.Init(logging.Config{
		AppName: "executive-api",
		Version: os.Getenv("APP_VERSION"),
		Env:     os.Getenv("APP_ENV"),
		Level:   os.Getenv("LOG_LEVEL"),
		Compact: true,
	})

	// Redis client from env
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	limiter := redis_rate.NewLimiter(rdb)

	router := chi.NewRouter()

	// ===== Middleware (must come first) =====
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Use(logging.RequestLogger(&httplog.Options{
		Level:           slog.LevelInfo,
		Schema:          httplog.SchemaECS,
		RecoverPanics:   true,
		Skip:            func(_ *http.Request, status int) bool { return status == 404 || status == 405 },
		LogRequestBody:  func(r *http.Request) bool { return r.Header.Get("Debug") == "reveal-body-logs" },
		LogResponseBody: func(r *http.Request) bool { return r.Header.Get("Debug") == "reveal-body-logs" },
	}))

	router.Use(func(next http.Handler) http.Handler {
		return logging.AddAttrs(next) // per-request attrs
	})

	router.Use(logging.RateLimitMiddleware(limiter, redis_rate.PerMinute(300), ClientKey))

	// ===== Test route =====
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"api is working well"}`))
	})

	// ===== Swagger =====
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// ===== API v1 routes =====
	router.Route("/api/v1", func(r chi.Router) {
		// stores
		userStore := user.NewStore(s.db)
		productStore := product.NewStore(s.db)
		categoryStore := category.NewStore(s.db)
		reviewStore := review.NewStore(s.db)
		cartStore := cart.NewStore(s.db)
		orderStore := order.NewStore(s.db)
		addressStore := address.NewStore(s.db)
		paymentStore := payment.NewStore(s.db)

		// handlers
		userHandler := user.NewHandler(userStore)
		productHandler := product.NewHandler(productStore)
		categoryHandler := category.NewHandler(categoryStore)
		reviewHandler := review.NewHandler(reviewStore, userStore)
		cartHandler := cart.NewHandler(cartStore, userStore)
		orderHandler := order.NewHandler(orderStore, userStore, addressStore)
		addressHandler := address.NewHandler(addressStore, userStore)
		paymentHandler := payment.NewHandler(paymentStore, orderStore)

		// per-request attrs for authenticated user
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
				if uid := types.UserIDFromContext(rr.Context()); uid != uuid.Nil {
					next = logging.AddAttrs(next, slog.String("user_id", uid.String()))
				}
				next.ServeHTTP(w, rr)
			})
		})

		// register routes
		userHandler.RegisterRoutes(r)
		productHandler.RegisterRoutes(r)
		categoryHandler.RegisterRoutes(r)
		reviewHandler.RegisterRoutes(r)
		cartHandler.RegisterRoutes(r)
		orderHandler.RegisterRoutes(r)
		addressHandler.RegisterRoutes(r)
		paymentHandler.RegisterRoutes(r)
	})

	log.Printf("Server listening on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
