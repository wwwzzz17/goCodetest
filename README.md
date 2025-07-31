# Product API Service

A simple REST API service for managing products built with Go and Gin framework.

## Features

- Create, Read, Update, Delete (CRUD) operations for products
- Product search with filters (name, price range, quantity range)
- In-memory storage with thread-safe operations
- JSON API responses
- Health check endpoint

## Prerequisites

- Go 1.19 or higher
- Git

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd goCodetest
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Start the Server

```bash
# From the project root directory
go run cmd/main.go
```

The server will start on `http://localhost:8080`

### 4. Verify Server is Running

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy"
  }
}
```

## API Endpoints

### Products

- `POST /products` - Create a new product
- `GET /products` - Get all products with optional filters
- `GET /products/:id` - Get a specific product by ID
- `PUT /products/:id` - Update a product by ID
- `DELETE /products/:id` - Delete a product by ID

### Health Check

- `GET /health` - Health check endpoint

## Testing

### Automated Testing

Run the provided test script to test all API functionality:

```bash
# Make the test script executable
chmod +x test_api.sh

# Run the tests
./test_api.sh
```

The test script will:
1. Create 4 sample products
2. Test search functionality with various filters
3. Update a product
4. Delete a product
5. Display the final product list

## Rate limiter
```bash
./internal/services/rate_limiter.go