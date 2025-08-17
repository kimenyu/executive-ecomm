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
- **Fully containerized deployment** with Docker and Docker Compose

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
│   ├── Dockerfile    # Node.js service containerization
│   ├── index.js      # Main Node.js application
│   └── package.json  # Node.js dependencies
├── Dockerfile        # Go backend containerization
├── docker-compose.yml # Multi-service orchestration
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites

- Docker and Docker Compose installed
- Git for cloning the repository

### 1. Clone the repository

```bash
git clone https://github.com/kimenyu/executive-ecomm.git
cd executive-ecomm
```

### 2. Setup Environment Variables

The application uses environment variables for configuration. For Docker deployment, the essential variables are already configured in `docker-compose.yml`, but you may want to create `.env` files for additional customization.

#### Optional: Go Backend `.env` (for local development)

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

#### Optional: Node.js Mpesa Service `.env` (for production Mpesa integration)

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
GO_BACKEND_NOTIFY_URL=http://goapp:8080/api/v1/payments/confirm
NODE_NOTIFY_SECRET=supersecret-node-key
JWT_SECRET=supersecretkey
```

### 3. Run with Docker (Recommended)

The application is fully containerized with Docker Compose, which will start all services including PostgreSQL, Redis, the Go backend, and the Node.js Mpesa service.

```bash
# Build and start all services
docker compose up --build

# Or run in detached mode
docker compose up --build -d

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Stop and remove volumes (clean reset)
docker compose down -v
```

#### Docker Services

The Docker setup includes:

- **PostgreSQL Database** (`db`): Accessible on port 5432
- **Redis Cache** (`redis`): Accessible on port 6379
- **Go Backend** (`goapp`): Accessible on port 8080
- **Node.js Mpesa Service** (`mpesa-service`): Accessible on port 5000

All services are connected via Docker networking and include health checks for reliability.

### 4. Alternative: Local Development Setup

If you prefer to run services locally:

**Start Go backend:**
```bash
make build
make run
```

**Start Node.js Mpesa service:**
```bash
cd mpesa-service
npm install
npm start
```

**Note**: You'll need to run PostgreSQL and Redis separately for local development.

### 5. View API Documentation

Once the services are running, visit:
```
http://localhost:8080/swagger/index.html
```

## Docker Configuration

### Container Architecture

The application uses a multi-stage Docker build:

- **Go Backend**: Multi-stage build using `golang:1.24-alpine` for building and `alpine:3.19` for runtime
- **Node.js Service**: Built on `node:20-alpine` for optimal size and performance
- **PostgreSQL**: Official `postgres:15` image with custom initialization
- **Redis**: Official `redis:7` image with persistence enabled

### Volume Management

Docker Compose creates persistent volumes for:
- `postgres_data`: Database storage
- `redis_data`: Redis persistence

### Health Checks

All services include health checks:
- **PostgreSQL**: `pg_isready` command
- **Redis**: `redis-cli ping` command
- **Service Dependencies**: Proper startup ordering with `depends_on`

### Container Networking

Services communicate via Docker's internal networking:
- Go backend connects to PostgreSQL at `db:5432`
- Go backend connects to Redis at `redis:6379`
- Node.js service connects to Go backend at `goapp:8080`

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

### Docker Development Commands

```bash
# Rebuild specific service
docker compose build goapp
docker compose build mpesa-service

# View service logs
docker compose logs goapp
docker compose logs mpesa-service

# Execute commands in running containers
docker compose exec goapp /bin/sh
docker compose exec db psql -U ecommerce_user -d ecommerce

# Scale services (if needed)
docker compose up --scale goapp=2
```

## Deployment

### Production Notes

1. **Build for production:**
   ```bash
   docker compose -f docker-compose.yml up --build -d
   ```

2. **Environment Variables**: Create production `.env` files with:
    - Production database credentials
    - Production Mpesa API credentials
    - Strong JWT secrets
    - Production Redis configuration

3. **SSL/TLS**: Configure reverse proxy (nginx/traefik) for HTTPS termination

4. **Monitoring**: Add container monitoring with tools like Prometheus/Grafana

### Production Considerations

1. **Environment Variables**: Use production Mpesa credentials
2. **HTTPS**: Ensure both Node.js and Go services use HTTPS
3. **Database**: Use production PostgreSQL with connection pooling
4. **Redis**: Use Redis for payment mapping storage in production(The application is currently using js map)
5. **Monitoring**: Implement payment monitoring and alerting
6. **Backup**: Regular database backups including payment records
7. **Container Orchestration**: Consider using Kubernetes for large-scale deployments
8. **Load Balancing**: Use multiple container replicas behind a load balancer
9. **Security**: Implement proper network policies and secrets management

### Cloud Deployment Options

- **Docker Swarm**: For simple multi-node deployments
- **Kubernetes**: For advanced orchestration and scaling
- **Cloud Container Services**: AWS ECS/Fargate, Google Cloud Run, Azure Container Instances
- **Platform-as-a-Service**: Railway, Render, DigitalOcean App Platform

## Contributing

Pull requests and issues are welcome. Before contributing, ensure your code is tested and follows existing patterns.


## Troubleshooting

### Common Issues

1. **Container Startup Issues**: Check logs with `docker compose logs <service_name>`

2. **Database Connection**: Ensure PostgreSQL container is healthy before Go app starts

3. **Port Conflicts**: Change ports in `docker-compose.yml` if 5432, 6379, 8080, or 5000 are in use

4. **Missing AccountReference**: Mpesa sometimes doesn't return `AccountReference` in callbacks. The system handles this with internal mapping.

5. **JWT Token Issues**: Ensure `JWT_SECRET` matches between Node.js and Go services.

6. **Callback Timeouts**: Mpesa callbacks have short timeouts. Process quickly and return success immediately.

7. **Amount Mismatches**: System validates payment amounts against order totals.

### Docker Troubleshooting

```bash
# Check container status
docker compose ps

# View detailed logs
docker compose logs --tail=100 -f goapp

# Restart specific service
docker compose restart goapp

# Clean rebuild
docker compose down -v
docker compose up --build
```

## License

This project is open-source and available under the [MIT License](LICENSE).

## Maintainer

Created and maintained by [@kimenyu](https://github.com/kimenyu)