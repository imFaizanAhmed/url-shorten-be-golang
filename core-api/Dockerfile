# Multi-stage build for Go application
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o core-api .
# RUN chmod +x core-api

# Final stage - minimal image
FROM alpine:latest

# # Install ca-certificates for HTTPS requests
# RUN apk --no-cache add ca-certificates

# # Create a non-root user
# RUN adduser -D -s /bin/sh appuser

# # Create app directory and set permissions
# RUN mkdir -p /app && chown -R appuser:appuser /app

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/core-api .

# # Make binary executable and change ownership
# RUN chmod +x core-api && chown appuser:appuser core-api

# # Switch to non-root user
# USER appuser

# Expose port 8080 (Beego's default port)
EXPOSE 8080

# Run the application
CMD ["./core-api"]
