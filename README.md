# ShopCx FrontEnd Service

A Go-based frontend service for the ShopCx demo application with intentionally vulnerable endpoints. This service handles user authentication, product management, file uploads, and serves as the main web interface for the e-commerce platform.

## Overview

The FrontEnd Service is a Go application built with the Gin web framework that provides the primary web interface and API gateway for the ShopCx platform. It handles user authentication, product operations, file uploads, and proxies requests to other microservices.

## Key Features

- **User Authentication**: JWT-based login and session management
- **Product Management**: Create, search, and manage product listings
- **File Upload**: Handle product images and document uploads
- **Search Functionality**: Product search with dynamic results
- **Admin Operations**: Administrative user and product management
- **Rate Limiting**: Basic request rate limiting and throttling
- **CORS Support**: Cross-origin resource sharing configuration

## Technology Stack

- **Go 1.21**: Core programming language
- **Gin Framework**: High-performance HTTP web framework
- **MySQL**: Database for user and product data
- **JWT**: JSON Web Token authentication
- **Gorilla Sessions**: Session management
- **HTML Templates**: Go template engine for web pages

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/login` | User authentication |
| GET | `/api/users/{id}` | Get user profile |
| POST | `/api/products` | Create new product |
| POST | `/api/upload` | Upload files |
| GET | `/search` | Product search page |
| POST | `/admin/delete-user` | Delete user (admin) |

## Dependencies

### Required Services
- **MySQL Database**: Required for user and product data
  - Default connection: `localhost:3306`
  - Database: `shopcx`
  - User: `admin` / Password: `admin123`

### Go Dependencies
See `go.mod` and `go.sum` for full dependency list including:
- Gin web framework
- MySQL driver
- JWT libraries
- Session management

## Build & Run

### Prerequisites
- Go 1.21 or higher
- MySQL server running locally
- Database `shopcx` created with appropriate tables

### Local Development
```bash
# Download dependencies
go mod download

# Build the application
go build -o main .

# Run the service
./main
```

The service will start on `http://localhost:8080`.

### Environment Variables
```bash
export PORT=8080
export DB_USER=admin
export DB_PASSWORD=admin123
export JWT_SECRET=very_secret_key_123
```

### Docker
```bash
# Build Docker image
docker build -t shopcx-frontend .

# Run container
docker run -p 8080:8080 \
  -e DB_USER=admin \
  -e DB_PASSWORD=admin123 \
  shopcx-frontend
```

## Configuration

### Database Setup
Create the required MySQL tables:
```sql
CREATE DATABASE shopcx;
USE shopcx;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    description TEXT
);
```

### File Upload Directory
The service expects an `uploads/` directory for file storage:
```bash
mkdir uploads
```

### Template Files
HTML templates are stored in the `templates/` directory:
- `search.html`: Product search results page

## Authentication

The service uses JWT tokens for authentication:
- **Login Endpoint**: `/login`
- **Token Format**: HS256 signed JWT
- **Secret**: Configurable via `JWT_SECRET` environment variable

### Login Example
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin&password=admin123"
```

## File Upload

Supports file uploads to the `/api/upload` endpoint:
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@image.jpg"
```

## Health Check

The service includes a health check endpoint:
- **Endpoint**: `/health`
- **Returns**: Service status and database connectivity

## Rate Limiting

Basic rate limiting is implemented:
- **Limit**: 100 requests per IP per minute
- **Response**: 429 Too Many Requests when exceeded

## Search Functionality

Product search is available via:
- **Web Interface**: `/search?q=query`
- **Returns**: HTML page with search results
- **Template**: Uses `search.html` template

## Security Note

⚠️ **This is an intentionally vulnerable application for security testing purposes. Do not deploy in production environments.**

### Known Vulnerabilities (Intentional)
- SQL injection in login endpoint
- Insecure direct object references
- Missing input validation
- Hardcoded secrets
- Excessive error information disclosure

## Recommended Checkmarx One Configuration
- Criticality: 4
- Cloud Insights: Yes
- Internet-facing: Yes
