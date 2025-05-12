# Web/API Gateway Connectivity Fix Plan

## Analysis of the Issue

After examining the code and configuration, I've identified the root cause of the connection issues between the web frontend and the API gateway:

### Issue 1: Server-Side Rendering (SSR) Fetch Calls in Docker Context

The primary issue is with SvelteKit's server-side rendering. When the app is rendering on the server within a Docker container, requests from the server to the API gateway are failing with ECONNREFUSED.

In `web/src/routes/+layout.server.ts`, SvelteKit is making server-side requests via:
```typescript
const airlinesResp = await apiClient.GET('/airlines', { fetch })
```

### Issue 2: Environment Variable Configuration 

The API client in `web/src/lib/api.ts` is configured with:
```typescript
export const apiClient = createClient<paths>({
	baseUrl: env.API_GATEWAY_URL || `http://localhost:${env.PUBLIC_API_PORT ?? '8000'}`,
})
```

While we've set the environment variable `API_GATEWAY_URL` to "http://api-gateway:8000", which is correct for Docker container-to-container communication, there's still an issue with how the endpoints are being handled.

### Issue 3: Path Configuration Mismatch

The web client is making requests to paths like `/airlines` and `/aircraft-types` (without the `/api` prefix), but the API gateway routes are configured with the `/api` prefix:

```go
services := map[string]string{
  "/api/airlines":       getServiceURL("AIRLINE_SERVICE_URL", "http://airline-service:8001"),
  "/api/aircraft-types": getServiceURL("AIRLINE_SERVICE_URL", "http://airline-service:8001") + "/aircraft-types",
  // ...
}
```

## Updated Diagnosis

After extensive testing, we've identified complex networking issues:

1. The API gateway is accessible from the host machine at http://localhost:8000/api/airlines
2. The API gateway is ping-able from the web container (via hostname 'api-gateway')
3. However, HTTP requests from the web container to the API gateway are timing out with ECONNREFUSED errors
4. **Network inspection shows the API gateway is listening only on IPv6 interfaces (:::8000)**
5. **Even after modifying the API gateway to explicitly listen on IPv4, connection issues persist**

This suggests a deeper Docker networking issue:

1. While the containers can resolve each other by hostname (ping works), HTTP connections are not getting through
2. The issue might involve:
   - Docker's internal network configuration or isolation
   - Container networking modes (bridge vs host)
   - Service discovery or DNS resolution problems
   - Port mapping or forwarding issues between containers

## Final Solution Implementation

After extensive debugging, we've implemented a combination of fixes that address both the path mismatch and networking issues:

1. **API Path Configuration**: Fixed in `web/src/lib/api.ts` by using the proper `path` parameter
   ```typescript
   export const apiClient = createClient<paths>({
     baseUrl: (env.API_GATEWAY_URL || `http://localhost:${env.PUBLIC_API_PORT ?? '8000'}`),
     path: env.PUBLIC_API_BASE_URL || '/api',
   })
   ```

2. **Client-Side Only Rendering**: Disabled server-side rendering in `web/src/routes/+layout.ts`
   ```typescript
   // Disable SSR to avoid server-side network issues
   export const ssr = false;
   ```

3. **URL Transformation for Browser**: Added logic to ensure client-side requests use localhost
   ```typescript
   initFetch: (url, init) => {
     // For client-side requests from the browser, ensure we use localhost
     if (typeof window !== 'undefined') {
       url = url.replace('http://api-gateway:8000', 'http://localhost:8000');
     }
     return fetch(url, init);
   },
   ```

4. **Docker Environment Variables**: Properly configured for both contexts
   ```yaml
   environment:
     API_GATEWAY_URL: "http://localhost:8000"
     PUBLIC_API_PORT: 8000
     PUBLIC_API_BASE_URL: "/api"
   ```

While we attempted to fix the network binding in the API gateway by using explicit IPv4 binding, deeper Docker networking issues remained. As a pragmatic solution, we've disabled server-side rendering and focused on ensuring the client-side experience works correctly.

## Solution Plan

### 1. Fix the API Path Mismatch

Two possible approaches:

#### Option A: Modify the Frontend API Client (Recommended)

Update the `apiClient` in `web/src/lib/api.ts` to properly handle API paths:

```typescript
export const apiClient = createClient<paths>({
	baseUrl: (env.API_GATEWAY_URL || `http://localhost:${env.PUBLIC_API_PORT ?? '8000'}`),
	path: env.PUBLIC_API_BASE_URL || '/api',

	// In server-side rendering context, we need the full URL to the API from the Docker container
	initFetch: (url, init) => {
		// During SSR, change the host to the API gateway service name for container communication
		if (typeof window === 'undefined') {
			// Use Docker service name instead of localhost for SSR requests
			url = url.replace(/http:\/\/localhost:8000|http:\/\/127.0.0.1:8000/, 'http://api-gateway:8000');
			console.log('SSR request URL:', url);
		} else {
			// For client-side requests from the browser, ensure we use localhost
			url = url.replace('http://api-gateway:8000', 'http://localhost:8000');
			console.log('Client request URL:', url);
		}
		return fetch(url, init);
	},
})
```

This approach:
1. Keeps the existing API gateway routes intact
2. Adds the proper `/api` path prefix using the environment variable
3. Handles URL translation between server-side and client-side rendering contexts

### 2. Fix SSR Network Issues

To address the server-side rendering connection issues, we implemented two main solutions:

#### Option A: Disable SSR for Initial Testing

Disabled server-side rendering to first ensure client-side rendering works:

```typescript
// web/src/routes/+layout.ts
export const ssr = false;
```

This approach allows us to focus on getting client-side communication working correctly first.

#### Option B: URL Transformation in API Client

Implemented a URL transformation mechanism in the `initFetch` function of the API client that:

1. In server context (SSR): Transforms `localhost` URLs to Docker service names (`http://api-gateway:8000`)
2. In client context: Ensures URLs use `localhost` for the browser to connect

### 3. Environment Variable Configuration

Updated environment variables in `docker-compose.yml` to support the dual-mode operation:

```yaml
web:
  # ... other settings
  environment:
    NODE_ENV: development
    # For server-side requests inside Docker network
    # But client browser will still use localhost
    API_GATEWAY_URL: "http://api-gateway:8000"
    PUBLIC_API_PORT: 8000
    PUBLIC_API_BASE_URL: "/api"
    PUBLIC_API_GATEWAY: "http://api-gateway:8000"
```

### 4. Implementation Summary

The implementation included these key components:

1. Added proper path prefix handling in the API client with `path: env.PUBLIC_API_BASE_URL || '/api'`
2. Added conditional URL transformation in `initFetch` to handle both server-side and client-side requests
3. Updated environment variables in Docker compose to support the dual routing needs
4. Disabled server-side rendering temporarily to isolate and fix client-side communication
5. Set up client-side data fetching in `+page.js` for improved reliability

### 5. Testing Steps

For testing and verifying the fix:

1. Rebuild the web container:
   ```bash
   docker-compose build web
   ```

2. Restart the web and API gateway services:
   ```bash
   docker-compose up -d web api-gateway
   ```

3. Check logs for connection errors:
   ```bash
   docker-compose logs web
   ```

4. Verify API Gateway endpoints:
   ```bash
   curl http://localhost:8000/api/airlines
   ```

5. Test the web interface at http://localhost:5179

### 6. Future Improvements

While our current solution gets the application working immediately with client-side rendering, several improvements should be considered for a production environment:

1. **Investigate Proper Docker Network Configuration**: Work with DevOps to ensure proper configuration of Docker networks, possibly using Docker's `host.docker.internal` DNS name or custom network configuration.

2. **Enable Server-Side Rendering**: Once network issues are resolved, re-enable SSR for better performance and SEO benefits.

3. **Implement Service Health Checks**: Add proper health checks for all services to ensure connectivity issues are detected early.

4. **Add Circuit Breaking & Fallbacks**: Implement circuit breaking and fallback mechanisms to handle service outages gracefully.

5. **Consider API Gateway Alternatives**: Evaluate more robust API gateway solutions like Kong, Nginx, or a full service mesh solution.

### 7. Advanced Docker Network Debugging

If the above solutions don't work, we'll need to perform more in-depth network debugging:

1. Verify API gateway container listening interfaces:
   ```bash
   docker-compose exec api-gateway netstat -tlnp
   ```

2. Check if any firewall rules are blocking inter-container communication:
   ```bash
   docker-compose exec api-gateway iptables -L
   ```

3. **Fix the API gateway IPv4 listening issue** (primary solution):
   ```go
   // In api-gateway/main.go, line ~51, modify the listening address to specify IPv4
   // Change from:
   // addr := ":" + port
   // To:
   addr := "0.0.0.0:" + port
   log.Printf("API Gateway listening on %s", addr)
   if err := http.ListenAndServe(addr, r); err != nil {
       log.Fatalf("Failed to start server: %v", err)
   }
   ```

4. Add explicit network configuration in docker-compose.yml:
   ```yaml
   networks:
     app-network:
       driver: bridge
   
   services:
     web:
       networks:
         - app-network
     api-gateway:
       networks:
         - app-network
   ```

5. **Create a special configuration for Docker network mode (recommended solution):**
   ```yaml
   # In docker-compose.yml
   web:
     # Add host network mode to bypass Docker network isolation
     network_mode: "host"
     environment:
       API_GATEWAY_URL: "http://localhost:8000"  # Use localhost when in host mode
       # Other settings...
   ```

6. As a workaround, bypass the API gateway during SSR:
   - Modify the layout.server.ts to provide mock data during server-side rendering
   - Use the API gateway only for client-side requests

7. For a production-ready solution, consider using a service mesh or dedicated API proxy:
   - Implement a proper service mesh like Istio or Linkerd
   - Use a dedicated API proxy like Nginx or Traefik that's proven to work well with Docker networks
   - Configure proper health checks and circuit breakers