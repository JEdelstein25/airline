#!/bin/bash

# Start all services using docker-compose
echo "Starting airline microservices..."

# Build and start all services in detached mode
docker-compose up -d

echo "Waiting for services to start..."
sleep 5

# Check health of services
echo "Checking health of services..."
curl http://localhost:8000/services/health

echo ""
echo "Services are running!"
echo "API Gateway: http://localhost:8000"
echo "Airline Service: http://localhost:8001"
echo "Airport Service: http://localhost:8002"
echo "Fleet Service: http://localhost:8003"
echo "Schedule Service: http://localhost:8004"
echo "Booking Service: http://localhost:8005"

echo ""
echo "To stop services: docker-compose down"