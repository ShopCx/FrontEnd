FROM golang:1.16-alpine

USER root

# Install curl for health checks
RUN apk add --no-cache gcc musl-dev curl

# Copy the application
WORKDIR /app
COPY . .

# Build the application
RUN go build -o main .

# Expose port
EXPOSE 8080

ENV DB_USER=admin
ENV DB_PASSWORD=admin123
ENV JWT_SECRET=very_secret_key_123
ENV ADMIN_PASSWORD=superadmin123

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./main"] 