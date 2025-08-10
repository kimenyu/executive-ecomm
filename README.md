# Executive eCommerce API

The **Executive eCommerce API** is a full-featured, backend written in Go for managing eCommerce workflows — including user authentication, product listings, cart management, order processing, product reviews, and Mpesa payment integration. Designed for scalability and maintainability, it features clean architecture, structured routing, and comprehensive API documentation via Swagger.

## Features

- **Distributed rate limiting** with Redis (`go-redis/redis_rate`) — configurable requests per minute per user/IP
- **Structured logging** using `log/slog` and `go-chi/httplog` with ECS format
- **Environment-based config** for log level, compact logs, and rate limits

- User registration and authentication with JWT
- Product and category management
- Cart creation and item tracking
- Order placement and tracking
- Product reviews with ownership validation
- **Mpesa payment integration** with Node.js STK Push service and Go backend confirmation
- PostgreSQL database integration
- Full Swagger/OpenAPI documentation

## Tech Stack

- **Rate Limiting**: Redis with `go-redis/redis_rate`
- **Logging**: Structured JSON logs via `log/slog` and `go-chi/httplog`

- **Language**: Go (Golang)
- **Framework**: Chi Router
- **Database**: PostgreSQL
- **ORM/Query Layer**: `database/sql`
- **Auth**: JWT with middleware
- **Documentation**: Swagger (`swaggo/swag`)
- **Payment Service**: Node.js Mpesa STK Push and Callback handler
- **Dependency Management**: Go Modules
- **Containerization**: Docker & Docker Compose

## Folder Structure

```
executive-ecomm/
├── cmd/              # Application entry point
│   └── api/          # Server and routing setup
├── services/         # Modular domain logic (user, product, order, etc.)
├── types/            # Models and payload DTOs
├── db/               # Database connection setup
├── configs/          # Config utilities
├── utils/            # Helper utilities
├── docs/             # Swagger docs
├── mpesa-service/    # Node.js Mpesa STK Push and callback handler
├── docker-compose.yml
├── Makefile
```

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/kimenyu/executive-ecomm.git
cd executive-ecomm
```

### 2. Setup Environment

Create a `.env` file

```
# ===== API CONFIG =====
API_ADDR=:8080
APP_NAME=executive-api
APP_ENV=development
APP_VERSION=v1.0.0
LOG_LEVEL=info
LOG_COMPACT=true

# ===== DATABASE =====
DATABASE_URL=postgres://ecommerce_user:strongpassword@localhost:5432/ecommerce?sslmode=disable

# ===== JWT AUTH =====
JWT_SECRET=supersecretkey
JWT_EXPIRATION_SECONDS=86400

# ===== REDIS =====
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# ===== RATE LIMIT CONFIG =====
RATE_LIMIT_REQUESTS_PER_MIN=300

# ===== SWAGGER =====
SWAGGER_HOST=localhost:8080
SWAGGER_SCHEMES=http

# ===== MPESA PAYMENT SERVICE =====
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_SHORTCODE=174379
MPESA_PASSKEY=your_passkey
MPESA_ENV=sandbox
CALLBACK_BASE_URL=http://your-node-service-url
GO_BACKEND_URL=http://your-go-backend-url/payments/confirm
NODE_NOTIFY_SECRET=some-secure-secret

```

### 3. Build & Run

```bash
make build
make run
```

Start Mpesa Node.js payment service separately(in mpesa-service folder)

```bash
cd mpesa-service
npm install
npm run start
```

### 4. View API Docs

Visit:

```
http://localhost:8080/swagger/index.html
```

### Payment flow with Mpesa
The Node.js service handles Mpesa STK Push requests and callbacks

On payment confirmation, it securely notifies the Go backend with payment details

Go backend creates payment records and updates order status accordingly
## API Documentation

All endpoints are documented using Swagger (OpenAPI 2.0). JWT-protected routes require an `Authorization: Bearer <token>` header.

## Contributing

Pull requests and issues are welcome. Before contributing, ensure your code is tested and follows existing patterns.

## License

This project is open-source and available under the [MIT License](LICENSE).

## Maintainer

Created and maintained by [@kimenyu](https://github.com/kimenyu)
