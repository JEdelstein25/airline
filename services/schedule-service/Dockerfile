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
COPY services/schedule-service/ /app/

# Simplify build without go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o schedule-service .

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/schedule-service .

# Set environment variables
ENV PORT=8004

# Expose port
EXPOSE 8004

# Run the application
CMD ["./schedule-service"]