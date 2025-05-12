#!/bin/bash

set -e  # Exit on error

echo "=== Building all microservices with shared module support ==="

# Ensure we're in the root directory
ROOT_DIR=$(pwd)

echo "Building from root directory: $ROOT_DIR"

# Function to build a specific service
build_service() {
    service_name=$1
    service_path=$2
    dockerfile_path="$service_path/Dockerfile"
    
    echo "\n=== Building $service_name ==="
    
    # Create a temporary build context directory
    build_context="$(mktemp -d)"
    echo "Created temp build context: $build_context"
    
    # Copy ONLY the essential shared module components
    echo "Copying minimal shared module..."
    mkdir -p "$build_context/shared/db" "$build_context/shared/models" "$build_context/shared/util" "$build_context/shared/api"
    cp "$ROOT_DIR/shared/go.mod" "$build_context/shared/"
    cp "$ROOT_DIR/shared/db"/*.go "$build_context/shared/db/"
    cp "$ROOT_DIR/shared/models"/*.go "$build_context/shared/models/"
    cp "$ROOT_DIR/shared/util"/*.go "$build_context/shared/util/"
    cp "$ROOT_DIR/shared/api"/*.go "$build_context/shared/api/"
    
    # Copy the service files
    echo "Copying $service_name files..."
    mkdir -p "$build_context/$service_path"
    cp -R "$ROOT_DIR/$service_path"/* "$build_context/$service_path/"
    
    # Build the Docker image
    echo "Building Docker image for $service_name..."
    docker build -t "airline/$service_name" -f "$dockerfile_path" "$build_context"
    
    # Clean up the temporary directory
    rm -rf "$build_context"
    echo "Build for $service_name completed."
}

# Build the shared module first (not a Docker build, just ensure it's ready)
echo "\n=== Preparing shared module ==="
cd "$ROOT_DIR/shared"
go mod tidy || echo "Note: go mod tidy for shared module failed, but we'll continue"

# Build each service
build_service "api-gateway" "api-gateway"
build_service "airline-service" "services/airline-service"
build_service "airport-service" "services/airport-service"
build_service "fleet-service" "services/fleet-service"
build_service "schedule-service" "services/schedule-service"
build_service "booking-service" "services/booking-service"

echo "\n=== All services built successfully ==="

echo "\n=== Ready to run: docker-compose up -d ==="