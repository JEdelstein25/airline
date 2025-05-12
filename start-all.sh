#!/bin/bash

echo "=== Starting Airline Microservices Application ==="

# Start all services
docker-compose up -d

echo ""
echo "=== Services Started Successfully! ==="
echo ""
echo "API Gateway: http://localhost:8000"
echo "Frontend:    http://localhost:5179"
echo ""
echo "Individual Services:"
echo "- Airline Service:  http://localhost:8001"
echo "- Airport Service:  http://localhost:8002"
echo "- Fleet Service:    http://localhost:8003"
echo "- Schedule Service: http://localhost:8004"
echo "- Booking Service:  http://localhost:8005"
echo ""
echo "To stop all services: docker-compose down"