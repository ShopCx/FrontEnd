# ShopCx Frontend Service

A Go-based frontend service for the ShopCx demo application with intentionally vulnerable endpoints.

## Vulnerabilities

### 1. Insecure Deserialization
- **Location**: `main.go`
- **Description**: Uses unsafe YAML deserialization
- **Test Example**: 
  ```yaml
  !!js/function >
  function() { 
    require('child_process').exec('rm -rf /') 
  }
  ```

### 2. Template Injection
- **Location**: `main.go`
- **Description**: Unsafe template rendering with user input
- **Test Example**: 
  ```go
  {{.UserInput}}
  ```

### 3. Missing Authentication
- **Location**: `main.go`
- **Description**: No authentication checks on sensitive endpoints
- **Test Example**: Direct access to `/api/admin/*` endpoints

### 4. No Rate Limiting
- **Location**: `main.go`
- **Description**: No rate limiting on any endpoints
- **Test Example**: Rapid requests to any endpoint

### 5. Insecure CORS Configuration
- **Location**: `main.go`
- **Description**: Overly permissive CORS settings
- **Test Example**: Cross-origin requests from any domain

## Setup

1. Install Go 1.16 or later
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the service:
   ```bash
   go run main.go
   ```

## API Documentation

API documentation is available at `/swagger/index.html` when the service is running.

## Dependencies

- Go 1.16+
- gin-gonic/gin: v1.7.0
- swaggo/gin-swagger: v1.4.3
- go-yaml/yaml: v2.4.0

## Security Best Practices (What NOT to do)

1. Using unsafe YAML deserialization
2. Disabling authentication checks
3. Having no rate limiting
4. Using overly permissive CORS
5. Not validating user input in templates 