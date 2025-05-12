# Docker Build Plan for Airline Microservices

## Current Issues

### 1. Go Version Issues

- **Problem**: Original Dockerfiles specified Go 1.23, which doesn't exist in Docker Hub
- **Solution**: Downgraded to Go 1.21 in all Dockerfiles

### 2. Module Dependency Issues

- **Problem**: Several services have dependencies on a shared module (`github.com/stellora/airline/shared` or `github.com/JEdelstein25/airline/shared`) with a replace directive pointing to `../../shared`
- **Issue**: This works in local development but fails in Docker builds as the replaced path doesn't exist in the Docker context
- **Services affected**:
  - `airline-service`
  - `booking-service`
  - `fleet-service`

### 3. Docker File Copy Issues

- **Problem**: Commands like `COPY go.mod go.sum* ./` were failing when go.sum doesn't exist
- **Solution**: Simplified to only copy the go.mod file which is required

### 4. go mod download Failures

- **Problem**: The `RUN go mod download` step fails due to the shared module replacement issues
- **Temporary solution**: Removed this step to allow builds to complete

## Solution Strategy

### Short-term Fixes (Implemented)

1. **Fix Go Version**:
   - Use `golang:1.21-alpine` instead of `golang:1.23-alpine` in all Dockerfiles

2. **Simplify Dockerfile Commands**:
   - Remove problematic copy commands for go.sum files
   - Remove the go mod download step temporarily

3. **Remove Shared Dependencies**:
   - Remove shared module dependencies from go.mod files for containerized builds

### Long-term Solutions

1. **Proper Shared Module Structure**:
   - Create a proper shared module with versioning
   - Publish it as a proper Go module or include it in the Docker build context
   - Options:
     - Use a proper Go module with versioning
     - Use Git submodules
     - Include shared code in each service if isolation is preferred

2. **Multi-stage Docker Builds**:
   - Optimize builds with more sophisticated multi-stage Dockerfiles
   - Add proper vendoring for dependencies

3. **Containerization Best Practices**:
   - Implement proper container health checks
   - Configure appropriate resource limits
   - Add proper logging and monitoring

## Implementation Plan

### Phase 1: Make All Services Buildable (Current)

- Simplify Dockerfiles to ensure basic functionality
- Remove problematic dependencies temporarily
- Ensure all services can be built and started

### Phase 2: Implement Proper Shared Code Strategy

- **Option A: Go Module Approach**
  - Create proper Go module structure for shared code
  - Set up versioning for the shared module
  - Update service go.mod files to use proper versioned imports

- **Option B: Include in Docker Context**
  - Modify Dockerfiles to include shared code in build context
  - Example:
    ```Dockerfile
    FROM golang:1.21-alpine AS builder
    WORKDIR /app
    
    # Copy shared modules first
    COPY shared/ /app/shared/
    
    # Copy service code
    COPY services/my-service/ /app/
    
    # Build with proper module path
    RUN CGO_ENABLED=0 GOOS=linux go build -o my-service .
    ```

### Phase 3: Optimize for Production

- Implement proper connection pooling for databases
- Add robust error handling and retry logic
- Implement proper logging and monitoring
- Set up CI/CD pipeline for automated builds