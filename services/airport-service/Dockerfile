FROM golang:1.21-alpine AS builder

WORKDIR /app

# Manually handle shared module (more reliable)
RUN mkdir -p /app/shared/db /app/shared/models /app/shared/api /app/shared/util
COPY shared/go.mod /app/shared/
COPY shared/db/*.go /app/shared/db/
COPY shared/models/*.go /app/shared/models/
COPY shared/api/*.go /app/shared/api/
COPY shared/util/*.go /app/shared/util/

# Copy the service code and go.mod
COPY services/airport-service/ /app/

# Simplify build without go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o airport-service .

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/airport-service .

# Set environment variables
ENV PORT=8002

# Expose port
EXPOSE 8002

# Run the application
CMD ["./airport-service"]