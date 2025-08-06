# Executive eCommerce API

The **Executive eCommerce API** is a full-featured, production-grade backend written in Go for managing eCommerce workflows — including user authentication, product listings, cart management, order processing, and product reviews. Designed for scalability and maintainability, it features clean architecture, structured routing, and comprehensive API documentation via Swagger.

## Features

- User registration and authentication with JWT
- Product and category management
- Cart creation and item tracking
- Order placement and tracking
- Product reviews with ownership validation
- PostgreSQL database integration
- Full Swagger/OpenAPI documentation

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Chi Router
- **Database**: PostgreSQL
- **ORM/Query Layer**: `database/sql`
- **Auth**: JWT with middleware
- **Documentation**: Swagger (`swaggo/swag`)
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

Create a `.env` file:

```
API_ADDR=localhost:8080
DB_URL=postgres://youruser:yourpass@localhost:5432/executive?sslmode=disable
JWT_SECRET=your-secret
```

### 3. Build & Run

```bash
make build
make run
```

Or, to just generate Swagger docs:

```bash
make swagger
```

### 4. View API Docs

Visit:

```
http://localhost:8080/swagger/index.html
```

## API Documentation

All endpoints are documented using Swagger (OpenAPI 2.0). JWT-protected routes require an `Authorization: Bearer <token>` header.

## Contributing

Pull requests and issues are welcome. Before contributing, ensure your code is tested and follows existing patterns.

## License

This project is open-source and available under the [MIT License](LICENSE).

## Maintainer

Created and maintained by [@kimenyu](https://github.com/kimenyu)
