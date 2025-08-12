# Executive eCommerce API

The **Executive eCommerce API** is a full-featured, backend written in Go for managing eCommerce workflows — including user authentication, product listings, cart management, order processing, product reviews, and Mpesa payment integration. Designed for scalability and maintainability, it features clean architecture, structured routing, and comprehensive API documentation via Swagger.

## Features

- **Distributed rate limiting** with Redis (`go-redis/redis_rate`) — configurable requests per minute per user/IP
- **Structured logging** using `log/slog` and `go-chi/httplog` with ECS format
- **Environment-based config** for log level, compact logs, and rate limits
- **User registration and authentication** with JWT
- **Product and category management**
- **Cart creation and item tracking**
- **Order placement and tracking**
- **Product reviews** with ownership validation
- **Complete Mpesa payment integration** with STK Push, callback handling, and payment confirmation
- **PostgreSQL database integration** with comprehensive payment tracking
- **Full Swagger/OpenAPI documentation**

## Tech Stack

- **Rate Limiting**: Redis with `go-redis/redis_rate`
- **Logging**: Structured JSON logs via `log/slog` and `go-chi/httplog`
- **Language**: Go (Golang)
- **Framework**: Chi Router
- **Database**: PostgreSQL
- **ORM/Query Layer**: `database/sql`
- **Auth**: JWT with middleware
- **Documentation**: Swagger (`swaggo/swag`)
- **Payment Service**: Node.js Mpesa STK Push with callback handler
- **Payment Processing**: Go backend with payment confirmation and database persistence
- **Dependency Management**: Go Modules
- **Containerization**: Docker & Docker Compose

## Folder Structure

```
executive-ecomm/
├── cmd/              # Application entry point
│   └── api/          # Server and routing setup
├── services/         # Modular domain logic (user, product, order, payment, etc.)
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

Create a `.env` file for the Go backend:

```bash
# ===== API CONFIG =====
API_ADDR=:8080
APP_NAME=executive-api
APP_ENV=development
APP_VERSION=v1.0.0
LOG_LEVEL=info
LOG_COMPACT=true
LOG_TIMEZONE=Africa/Nairobi

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
RATE_LIMIT_BURST=20

# ===== SWAGGER =====
SWAGGER_HOST=localhost:8080
SWAGGER_SCHEMES=http

# ===== PAYMENT CONFIRMATION =====
NODE_NOTIFY_SECRET=supersecret-node-key
```

Create a `.env` file for the Node.js Mpesa service:

```bash
# ===== NODE.JS MPESA SERVICE =====
PORT=5000

# Safaricom Daraja API credentials
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_SHORTCODE=174379
MPESA_PASSKEY=your_passkey
MPESA_ENV=sandbox    # or "production"

# Callback URL (make sure it's deployed)
CALLBACK_BASE_URL=https://your-node-service.onrender.com

# Go backend notification settings
GO_BACKEND_NOTIFY_URL=https://your-go-backend.ngrok.io/api/v1/payments/confirm
NODE_NOTIFY_SECRET=supersecret-node-key
JWT_SECRET=supersecretkey
```

### 3. Build & Run

Start the Go backend:
```bash
make build
make run
```

Start the Node.js Mpesa service(good choice if you deploy deplooy the backend. Daraja callbacks does not work on locahosts):
```bash
cd mpesa-service
npm install
npm start
```

### 4. View API Docs

Visit:
```
http://localhost:8080/swagger/index.html
```

## Payment Integration

### Mpesa STK Push Flow

The system implements a complete Mpesa payment workflow:

1. **Payment Initiation**: Client requests payment via Node.js service
   ```bash
   POST /mpesa/stkpush
   {
     "order_id": "uuid-here",
     "amount": 100,
     "phone": "254712345678"
   }
   ```

2. **STK Push**: Node.js service initiates Mpesa STK Push
    - Stores `CheckoutRequestID` → `order_id` mapping
    - Customer receives payment prompt on phone

3. **Callback Processing**: Mpesa sends callback to Node.js service
    - Maps `CheckoutRequestID` back to original `order_id`
    - Extracts payment details (receipt, phone, amount)

4. **Payment Confirmation**: Node.js notifies Go backend
   ```bash
   POST /api/v1/payments/confirm
   {
     "order_id": "uuid-here",
     "status": "success",
     "amount": 100,
     "provider": "mpesa",
     "mpesa_receipt": "THC9YIBI4J",
     "phone": "254712345678"
   }
   ```

5. **Database Storage**: Go backend stores payment record and updates order status

### Key Features

- **Order Tracking**: Maps Mpesa `CheckoutRequestID` to your internal `order_id`
- **Automatic Cleanup**: Removes old payment mappings to prevent memory leaks
- **Comprehensive Logging**: Full payment flow logging for debugging
- **Error Handling**: Graceful handling of failed payments and network issues
- **Security**: Validates payment amounts and authenticates callback requests

### Database Schema

The system creates payment records with full Mpesa metadata:

```sql
payments:
- id (uuid)
- order_id (uuid) → links to orders table
- amount (decimal)
- status (success/failed)
- provider (mpesa)
- checkout_request_id (Mpesa identifier)
- merchant_request_id (Mpesa identifier)  
- mpesa_receipt (Mpesa receipt number)
- phone (customer phone number)
- metadata (full Mpesa callback JSON)
- created_at (timestamp)
```

### Finding User Transactions

To find which user made a payment:
```sql
SELECT p.*, o.user_id, u.email 
FROM payments p
JOIN orders o ON p.order_id = o.id  
JOIN users u ON o.user_id = u.id
WHERE p.mpesa_receipt = 'THC9YIBI4J';
```

## API Documentation

All endpoints are documented using Swagger (OpenAPI 2.0). JWT-protected routes require an `Authorization: Bearer <token>` header.

### Main Endpoints

- **Authentication**: `/api/v1/auth/*`
- **Products**: `/api/v1/products/*`
- **Orders**: `/api/v1/orders/*`
- **Cart**: `/api/v1/cart/*`
- **Payments**: `/api/v1/payments/*`
- **Reviews**: `/api/v1/reviews/*`

### Payment Endpoints

- **Node.js Service**:
    - `POST /mpesa/stkpush` - Initiate STK Push
    - `POST /mpesa/callback` - Handle Mpesa callbacks
    - `GET /mpesa/mappings` - Debug payment mappings

- **Go Backend**:
    - `POST /api/v1/payments/confirm` - Receive payment confirmations

## Development

### Running Tests
```bash
make test
```

### Database Migrations
```bash
make migrate-up
make migrate-down
```

### Generating Swagger Docs
```bash
make swagger
```

## Deployment

### Production Considerations

1. **Environment Variables**: Use production Mpesa credentials
2. **HTTPS**: Ensure both Node.js and Go services use HTTPS
3. **Database**: Use production PostgreSQL with connection pooling
4. **Redis**: Use Redis for payment mapping storage in production
5. **Monitoring**: Implement payment monitoring and alerting
6. **Backup**: Regular database backups including payment records

### Docker Deployment
```bash
docker-compose up -d
```

## Contributing

Pull requests and issues are welcome. Before contributing, ensure your code is tested and follows existing patterns.

### Adding New Payment Methods

The payment system is designed to be extensible. To add new payment providers:

1. Add provider to `types.Payment.Provider` enum
2. Implement provider-specific handler in Node.js service
3. Update Go backend confirmation handler
4. Add appropriate database fields if needed

## Troubleshooting

### Common Issues

1. **Missing AccountReference**: Mpesa sometimes doesn't return `AccountReference` in callbacks. The system handles this with internal mapping.

2. **JWT Token Issues**: Ensure `JWT_SECRET` matches between Node.js and Go services.

3. **Callback Timeouts**: Mpesa callbacks have short timeouts. Process quickly and return success immediately.

4. **Amount Mismatches**: System validates payment amounts against order totals.

## License

This project is open-source and available under the [MIT License](LICENSE).

## Maintainer

Created and maintained by [@kimenyu](https://github.com/kimenyu)