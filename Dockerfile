# Using an outdated base image with known vulnerabilities
FROM golang:1.16-alpine

# Running as root (intentionally insecure)
USER root

# Copy the application
WORKDIR /app
COPY . .

# Using outdated package manager
RUN apk add --no-cache gcc musl-dev

# Build the application
RUN go build -o main .

# Expose port
EXPOSE 8080

# Hardcoded credentials in environment variables (intentionally insecure)
ENV DB_USER=admin
ENV DB_PASSWORD=admin123
ENV JWT_SECRET=very_secret_key_123
ENV ADMIN_PASSWORD=superadmin123

# Running as root (intentionally insecure)
CMD ["./main"] 