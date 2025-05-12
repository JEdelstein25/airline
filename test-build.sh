#!/bin/bash

set -e  # Exit on error

echo "=== Testing build for Docker images ==="

# Build a simplified base image for all services
cat > Dockerfile.simple << 'EOF'
FROM golang:1.21-alpine
WORKDIR /app

# Copy the shared module and go.mod for testing
COPY shared/ /app/shared/

# Print debug info about shared module
RUN ls -la /app/shared && ls -la /app/shared/db

# This is just a test - we're not actually trying to build anything
CMD ["echo", "Test completed"]
EOF

# Try to build the test container
docker build -t test-airline-build -f Dockerfile.simple .

echo "\n=== Success! The shared module paths are correct ==="

# Clean up
rm Dockerfile.simple