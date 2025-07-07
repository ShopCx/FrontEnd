FROM golang:1.16-alpine

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

ENV DB_USER=admin
ENV DB_PASSWORD=admin123
ENV JWT_SECRET=very_secret_key_123
ENV ADMIN_PASSWORD=superadmin123

CMD ["./main"] 