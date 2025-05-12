FROM golang:1.21-alpine

WORKDIR /app

# Copy app files
COPY . .

# Build
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /api-gateway

# Set environment variables
ENV PORT=8000
ENV AIRLINE_SERVICE_URL=http://airline-service:8001
ENV AIRPORT_SERVICE_URL=http://airport-service:8002
ENV FLEET_SERVICE_URL=http://fleet-service:8003
ENV SCHEDULE_SERVICE_URL=http://schedule-service:8004
ENV BOOKING_SERVICE_URL=http://booking-service:8005

# Expose port
EXPOSE 8000

# Run the application
CMD ["/api-gateway"]