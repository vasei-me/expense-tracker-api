# Expense Tracker API

A robust and scalable Expense Tracker REST API built with Go and SOLID architecture principles. This API allows users to track their expenses with filtering capabilities, user authentication, and full CRUD operations.

## 🚀 Features

### Authentication & Authorization

- ✅ User registration with email and password
- ✅ User login with JWT token generation
- ✅ Protected endpoints using JWT authentication
- ✅ Secure password handling

### Expense Management

- ✅ Create new expenses with categories
- ✅ Retrieve all expenses with pagination
- ✅ Update existing expenses
- ✅ Delete expenses
- ✅ Filter expenses by date ranges
- ✅ Filter expenses by categories
- ✅ Calculate total expenses

### Categories

- 🛒 Groceries
- 🎬 Leisure
- 📱 Electronics
- ⚡ Utilities
- 👕 Clothing
- 🏥 Health
- 📦 Others

### Date Filters

- 📅 Past week
- 📅 Past month
- 📅 Last 3 months
- 📅 Custom date range

## 📋 Prerequisites

- Go 1.21+
- SQLite3 (or PostgreSQL for production)
- Git

## 🛠️ Installation

### 1. Clone the repository

```bash
git clone https://github.com/vasei-me/expense-tracker-api.git
cd expense-tracker-api

```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables (optional)

Create a `.env` file in the root directory:

```env
PORT=5000
JWT_SECRET=your-secret-key-change-in-production
DB_TYPE=sqlite
DB_NAME=expense_tracker.db
```

### 4. Run the application

```bash
go run cmd/api/main.go
```

The API will start at `http://localhost:5000`

## 📊 API Endpoints

### Public Endpoints

| Method | Endpoint             | Description       |
| ------ | -------------------- | ----------------- |
| GET    | `/`                  | Welcome message   |
| GET    | `/health`            | Health check      |
| GET    | `/api/test`          | Test endpoint     |
| POST   | `/api/auth/register` | Register new user |
| POST   | `/api/auth/login`    | Login user        |

### Protected Endpoints (Require JWT)

| Method | Endpoint             | Description        |
| ------ | -------------------- | ------------------ |
| POST   | `/api/expenses`      | Create new expense |
| GET    | `/api/expenses`      | Get all expenses   |
| PUT    | `/api/expenses/{id}` | Update expense     |
| DELETE | `/api/expenses/{id}` | Delete expense     |

## 🔧 API Usage Examples

### 1. Register a new user

```bash
curl -X POST http://localhost:5000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 3. Create an expense (Protected)

```bash
curl -X POST http://localhost:5000/api/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-jwt-token" \
  -d '{
    "amount": 75.50,
    "category": "groceries",
    "description": "Weekly grocery shopping",
    "date": "2024-01-20"
  }'
```

### 4. Get all expenses with filters

```bash
# Get all expenses
curl -H "Authorization: Bearer demo-jwt-token" \
  http://localhost:5000/api/expenses

# Get expenses from last week
curl -H "Authorization: Bearer demo-jwt-token" \
  "http://localhost:5000/api/expenses?period=week"

# Get groceries from last month
curl -H "Authorization: Bearer demo-jwt-token" \
  "http://localhost:5000/api/expenses?period=month&category=groceries"

# Get custom date range
curl -H "Authorization: Bearer demo-jwt-token" \
  "http://localhost:5000/api/expenses?start_date=2024-01-01&end_date=2024-01-31"
```

### 5. Update an expense

```bash
curl -X PUT http://localhost:5000/api/expenses/exp-123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-jwt-token" \
  -d '{
    "amount": 85.00,
    "description": "Updated grocery list with organic items"
  }'
```

### 6. Delete an expense

```bash
curl -X DELETE http://localhost:5000/api/expenses/exp-123 \
  -H "Authorization: Bearer demo-jwt-token"
```

## 🏗️ Project Structure

```
expense-tracker/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── domain/
│   │   ├── entities/           # Business entities
│   │   ├── valueobjects/       # Value objects
│   │   └── repositories/       # Repository interfaces
│   ├── application/
│   │   ├── services/           # Business logic
│   │   ├── dto/                # Data transfer objects
│   │   └── interfaces/         # Application interfaces
│   └── infrastructure/
│       ├── database/           # Database implementations
│       ├── jwt/               # JWT implementation
│       ├── repositories/      # Repository implementations
│       └── http/
│           ├── handlers/      # HTTP handlers
│           └── middleware/    # HTTP middleware
├── migrations/                 # Database migrations
├── go.mod                     # Go modules
├── go.sum                     # Go dependencies
└── expense_tracker.db         # SQLite database file
```

## 🔒 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer your-jwt-token
```

## 🗄️ Database

### SQLite (Default)

The application uses SQLite by default for simplicity. The database file `expense_tracker.db` is automatically created.

### PostgreSQL

To use PostgreSQL, update the configuration:

```env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=expense_tracker
DB_SSLMODE=disable
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/application/services
```

### API Testing with Postman

1. Import the Postman collection from `docs/postman_collection.json`
2. Set up environment variables in Postman
3. Run the collection tests

## 📦 Deployment

### Using Docker

```bash
# Build Docker image
docker build -t expense-tracker-api .

# Run container
docker run -p 5000:5000 --env-file .env expense-tracker-api
```

### Docker Compose

```bash
# Start with PostgreSQL
docker-compose up -d
```

## 🔍 Monitoring

- Health endpoint: `GET /health`
- Database connection check
- Server status monitoring

## 🛡️ Security

- Password hashing with bcrypt
- JWT token authentication
- SQL injection prevention
- Input validation
- CORS support

## 🔄 Environment Variables

| Variable    | Default            | Description                     |
| ----------- | ------------------ | ------------------------------- |
| PORT        | 5000               | Server port                     |
| JWT_SECRET  | (random)           | JWT secret key                  |
| DB_TYPE     | sqlite             | Database type (sqlite/postgres) |
| DB_NAME     | expense_tracker.db | Database name/file              |
| DB_HOST     | localhost          | Database host                   |
| DB_PORT     | 5432               | Database port                   |
| DB_USER     | expense_user       | Database user                   |
| DB_PASSWORD | expense_password   | Database password               |
| DB_SSLMODE  | disable            | SSL mode for PostgreSQL         |

Project URL: https://roadmap.sh/projects/expense-tracker-api
