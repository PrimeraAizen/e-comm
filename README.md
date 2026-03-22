# E-Commerce API with Recommendation System

A modern e-commerce REST API built with Go, PostgreSQL, JWT authentication, and collaborative filtering recommendation system.

## Features

- **User Authentication**: JWT-based auth with access and refresh tokens
- **Product Catalog**: Full CRUD with categories, search, filtering, and pagination
- **PostgreSQL with UUIDs**: Relational database with UUID primary keys, migrations via Goose
- **User Interactions**: Track product views, likes, and purchases (via orders)
- **Recommendation System**: Collaborative filtering with weighted user interactions (purchases 50%, likes 35%, views 15%)
- **User Profiles**: Separate profile management with personal information
- **Product Statistics**: Materialized view for fast aggregated metrics (views, likes, purchases, ratings)
- **Swagger Documentation**: Complete OpenAPI/Swagger UI at `/swagger/index.html`
- **CORS Support**: Pre-configured for frontend integration
- **Graceful Shutdown**: Proper cleanup of resources on SIGTERM/SIGINT
- **RESTful API**: Clean architecture with Gin framework
- **Docker Support**: Easy deployment with Docker Compose

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose (recommended)
- PostgreSQL 15+ (if not using Docker)

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd e-comm
```

### 2. Configure Environment

```bash
cp config/config.example.yaml config/config.yaml
# Edit config/config.yaml with your settings
```

### 3. Start PostgreSQL

```bash
make docker-up
```

### 4. Run the Application

```bash
make run
```

Migrations run automatically on startup.

API available at: `http://localhost:8080`
Swagger UI available at: `http://localhost:8080/swagger/index.html`

## Default Credentials

Seeded automatically by migrations:

| Role      | Email                    | Password    |
|-----------|--------------------------|-------------|
| Admin     | admin@example.com        | password123 |
| Moderator | moderator@example.com    | password123 |
| User      | user1@example.com        | password123 |
| User      | user2@example.com        | password123 |
| Student   | student@example.com      | password123 |
| Teacher   | teacher@example.com      | password123 |

To re-seed manually:

```bash
go run scripts/seed/main.go
```

## API Documentation

### Swagger UI

Visit `http://localhost:8080/swagger/index.html` for interactive API documentation.

**Authentication in Swagger:**
1. Click the "Authorize" button
2. Enter: `Bearer <your_access_token>`
3. Click "Authorize" then "Close"

### Authentication Endpoints

```bash
# Register
POST /api/v1/auth/register
{"email": "user@example.com", "password": "password123"}

# Login
POST /api/v1/auth/login
{"email": "user@example.com", "password": "password123"}

# Refresh Token
POST /api/v1/auth/refresh
{"refresh_token": "your-refresh-token"}
```

### Product Endpoints (require authentication)

```bash
# List products with filters and pagination
GET /api/v1/products?page=1&limit=20
GET /api/v1/products?search=laptop&category_id=<uuid>
GET /api/v1/products?min_price=100&max_price=1000&sort_by=price&sort_order=asc

# Get product details
GET /api/v1/products/:id

# Get product statistics (from materialized view)
GET /api/v1/products/:id/statistics

# Create product (Admin only)
POST /api/v1/products
{"name": "iPhone 15 Pro", "description": "...", "category_id": "<uuid>", "price": 999.99, "stock": 100}

# Update product (Admin only)
PUT /api/v1/products/:id

# Delete product (Admin only)
DELETE /api/v1/products/:id
```

### User Interaction Endpoints

```bash
POST   /api/v1/products/:id/view       # Record view
POST   /api/v1/products/:id/like       # Like product
DELETE /api/v1/products/:id/like       # Unlike product
GET    /api/v1/products/:id/liked      # Check if liked
POST   /api/v1/products/:id/purchase   # Purchase (creates order + order_item)
GET    /api/v1/products/:id/purchased  # Check if purchased
```

### Category Endpoints (require authentication)

```bash
GET    /api/v1/categories        # List all categories
GET    /api/v1/categories/:id    # Get category
POST   /api/v1/categories        # Create (Admin only)
PUT    /api/v1/categories/:id    # Update (Admin only)
DELETE /api/v1/categories/:id    # Delete (Admin only)
```

### Profile Endpoints

```bash
GET    /api/v1/profiles/me               # Get my profile
PUT    /api/v1/profiles/me               # Update profile
PUT    /api/v1/profiles/me/password      # Change password
DELETE /api/v1/profiles/me/account       # Delete account
```

### Recommendation Endpoints

```bash
GET /api/v1/profiles/me/recommendations  # Personalized recommendations
GET /api/v1/profiles/me/interactions     # Full interaction history
GET /api/v1/profiles/me/views            # Viewed products
GET /api/v1/profiles/me/likes            # Liked products
GET /api/v1/profiles/me/purchases        # Purchase history
GET /api/v1/profiles/me/similar          # Similar users
```

## Configuration

`config/config.yaml`:

```yaml
http:
  host: "0.0.0.0"
  port: "8080"

database:
  host: localhost
  port: "5432"
  database: ecommerce
  username: postgres
  password: postgres
  ssl_mode: disable
  max_conns: 25
  min_conns: 5

jwt:
  secret: "your-secret-key-change-in-production"
  access_token_duration: "15m"
  refresh_token_duration: "168h"

logger:
  level: info
  format: json
  service: e-comm
  version: "1.0.0"
  environment: development
```

## Project Structure

```
.
├── cmd/web/main.go                      # Entry point
├── config/
│   ├── config.go                        # Configuration loader
│   └── config.yaml                      # Configuration file
├── internal/
│   ├── domain/
│   │   ├── user.go                      # User model (UUID string ID)
│   │   ├── profile.go                   # Profile model
│   │   ├── product.go                   # Product, Category, ProductFilter models
│   │   ├── interaction.go               # View, Like, Purchase, Summary models
│   │   ├── recommendation.go            # Recommendation models
│   │   └── errors.go                    # Domain errors
│   ├── repository/
│   │   ├── repository.go                # Repository container
│   │   ├── userRepository.go            # User CRUD
│   │   ├── profile_repository.go        # Profile CRUD
│   │   ├── product_repository.go        # Product & Category CRUD + statistics
│   │   └── interaction_repository.go    # Views, likes, purchases
│   ├── service/
│   │   ├── service.go                   # Service container
│   │   └── authService.go              # JWT auth logic
│   ├── delivery/rest/v1/
│   │   ├── handlers.go                  # Route registration
│   │   └── auth_api.go                  # Auth endpoints
│   └── server/server.go                 # HTTP server
├── pkg/
│   └── adapter/
│       ├── postgres.go                  # pgxpool connection + squirrel builder
│       └── migration.go                 # Goose migration runner
├── migrations/postgres/
│   ├── 20250916160326_initial_migration.sql   # users, profiles, roles schema
│   ├── 20251107000000_ecommerce_tables.sql    # products, categories, interactions, orders
│   └── 20251107000001_seed_data.sql           # Default roles, users, products
├── scripts/seed/main.go                 # Manual re-seed script
├── docs/                                # Swagger generated docs
├── docker-compose.yml
├── Makefile
└── go.mod
```

## Database Schema

All primary keys are UUIDs (`uuid_generate_v4()`). Key tables:

| Table | Description |
|-------|-------------|
| `users` | User accounts |
| `profiles` | Extended user info (name, phone, address) |
| `roles` / `user_roles` | Role-based access control |
| `categories` | Hierarchical product categories |
| `products` | Product catalog with full-text search indexes |
| `user_product_views` | Product view tracking |
| `user_product_likes` | Unique likes per user/product |
| `orders` / `order_items` | Purchase history |
| `product_reviews` | Ratings and reviews |
| `cart_items` | Shopping cart |
| `product_statistics` | Materialized view: aggregated metrics |

Migrations are managed by [Goose](https://github.com/pressly/goose) and run automatically on startup.

## Makefile Commands

```bash
make run          # Run the application
make build        # Build binary
make clean        # Remove build artifacts
make docker-up    # Start PostgreSQL
make docker-down  # Stop Docker containers
```

## Testing

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login and capture token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.access_token')

# List products
curl http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN"

# Purchase a product
curl -X POST http://localhost:8080/api/v1/products/<uuid>/purchase \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 1}'

# Get recommendations
curl http://localhost:8080/api/v1/profiles/me/recommendations \
  -H "Authorization: Bearer $TOKEN"
```

## Docker Deployment

```bash
docker-compose up -d      # Start everything
docker-compose logs -f    # View logs
docker-compose down       # Stop everything
```

## Troubleshooting

### PostgreSQL connection failed

```bash
docker ps | grep postgres
make docker-up
docker logs ecommerce_postgres
```

### Port 8080 already in use

```bash
lsof -ti:8080 | xargs kill -9
```

### Refresh materialized view manually

The `product_statistics` view is refreshed via the API, but can also be refreshed directly:

```sql
REFRESH MATERIALIZED VIEW CONCURRENTLY product_statistics;
```

## Recommendation System

Collaborative filtering based on weighted user interactions:

- **Purchases**: 50% weight
- **Likes**: 35% weight
- **Views**: 15% weight

The algorithm finds users with similar interaction patterns (weighted Jaccard similarity) and recommends products those users engaged with that the current user hasn't seen yet.

## Tech Stack

- [Gin](https://github.com/gin-gonic/gin) — HTTP framework
- [pgx/v5](https://github.com/jackc/pgx) — PostgreSQL driver
- [squirrel](https://github.com/Masterminds/squirrel) — SQL query builder
- [Goose](https://github.com/pressly/goose) — Database migrations
- [golang-jwt](https://github.com/golang-jwt/jwt) — JWT tokens
- [swaggo](https://github.com/swaggo/swag) — Swagger docs
- [viper](https://github.com/spf13/viper) — Configuration

---

*Last Updated: March 2026*
