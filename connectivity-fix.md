# API Gateway and Web Connectivity Issues

## Current Status

1. The API Gateway has been fixed and is now correctly proxying requests to the backend services.
2. The web service is still having connectivity issues when trying to reach the API gateway.

## Identified Problems

1. **Network Naming**: In Docker Compose, services refer to each other by their service name, not by localhost.
2. **Environment Variable Configuration**: The web service has incorrect URL configuration.

## Solution Steps

1. **Fix the API Gateway URL in the web service config:**
   - Change `API_GATEWAY_URL` from "http://localhost:8000/api" to "http://api-gateway:8000"
   - Change `PUBLIC_API_GATEWAY` from "http://localhost:8000" to "http://api-gateway:8000"

2. **Update the Web Service API Client:**
   The API client should use the correct URL format to reach the API gateway.

3. **Debug Connection Failures:**
   - The logs show ECONNREFUSED errors, indicating network connectivity issues
   - The API Gateway is correctly responding to manual HTTP requests

## Required Changes

1. Ensure the web service is using the correct API URL: "http://api-gateway:8000" 

2. The `/api` prefix is handled by the API Gateway, so paths in the client code should not include it.

3. Check Docker networking to ensure all services are on the same network.

## Testing Plan

1. Verify the API Gateway responds to requests at: http://localhost:8000/api/airlines
2. Check if the web service is using the correct URL to connect to the API Gateway
3. Check for any additional CORS or network issues

## Notes on Troubleshooting

The connection issue appears to be in the SvelteKit's server-side rendering phase. The web app is trying to fetch data during SSR but cannot connect to the API Gateway.