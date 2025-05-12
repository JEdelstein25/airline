# Airline Application: Monolith to Microservices Migration Plan

## Current Architecture (Monolith)

```
./
├── api/
│   └── airline.openapi.yaml       # OpenAPI specification
├── api-server/                    # Monolithic Go backend
│   ├── api/                       # API definitions
│   ├── db/                        # Database access and models
│   ├── handler.go                 # Request handlers
│   ├── aircraft.go                # Aircraft management
│   ├── airlines.go                # Airline management
│   ├── airports.go                # Airport management
│   ├── flights.go                 # Flight management
│   ├── schedules.go               # Schedule management
│   ├── passengers.go              # Passenger management
│   ├── seat_assignments.go        # Seat assignment management
│   └── ...                        # Other business logic
└── web/                           # Frontend application
    └── src/                       # Frontend source code
```

## Target Architecture (Microservices)

```
./
├── api-gateway/                   # API Gateway service
│   ├── main.go                    # Entry point
│   ├── routes.go                  # Route definitions
│   └── middleware/                # Common middleware
├── services/
│   ├── fleet-service/             # Aircraft and fleet management
│   │   ├── api/                   # Service API definition
│   │   ├── db/                    # Service database access
│   │   ├── handlers/              # HTTP handlers
│   │   ├── models/                # Domain models
│   │   └── main.go                # Service entry point
│   ├── airline-service/           # Airline management
│   │   ├── [similar structure]    
│   ├── airport-service/           # Airport management
│   │   ├── [similar structure]
│   ├── flight-service/            # Flight operations
│   │   ├── [similar structure]
│   ├── schedule-service/          # Schedule management
│   │   ├── [similar structure]
│   ├── booking-service/           # Passenger and booking management
│   │   ├── [similar structure]
│   └── notification-service/      # New service for notifications
│       ├── [similar structure]
├── shared/                        # Shared libraries and utilities
│   ├── auth/                      # Authentication
│   ├── messaging/                 # Message broker integration
│   ├── logging/                   # Logging utilities
│   └── models/                    # Shared domain models
├── docker/                        # Docker configurations
│   ├── docker-compose.yml         # Development compose file
│   └── services/                  # Service-specific Dockerfiles
├── database/                      # Database migrations and schemas
│   ├── migrations/                # Migration scripts
│   └── schemas/                   # Database schemas per service
└── web/                           # Frontend application (unchanged)
    └── src/                       # Frontend source code
```

## Service Boundaries

1. **API Gateway**
   - Route requests to appropriate microservices
   - Handle authentication and common middleware
   - Implement rate limiting and request tracking

2. **Fleet Service**
   - Manage aircraft types and fleet information
   - Track aircraft status and maintenance
   - APIs: GET/POST/PUT/DELETE /aircraft, /aircraft-types

3. **Airline Service**
   - Manage airline information and settings
   - APIs: GET/POST/PUT/DELETE /airlines

4. **Airport Service**
   - Manage airport information and timezone data
   - APIs: GET/POST/PUT/DELETE /airports

5. **Flight Service**
   - Manage flight operations
   - Track flight status and routing
   - APIs: GET/POST/PUT/DELETE /flights

6. **Schedule Service**
   - Manage flight schedules
   - Handle recurring flights and scheduling rules
   - APIs: GET/POST/PUT/DELETE /schedules

7. **Booking Service**
   - Manage passenger information
   - Handle seat assignments and reservations
   - APIs: GET/POST/PUT/DELETE /passengers, /seat-assignments, /itineraries

## Database Strategy

### Shared Schema vs Independent Schema Tradeoff

#### Initial Phase: Database per Service Pattern with Shared Schema
   - Each service connects to its own database instance but shares schema initially
   - Implement read replicas for high-traffic services
   - **Shared Schema Structure**:
     ```
     ┌─────────────────────────────────────────────────────────┐
     │                    Shared Schema                        │
     ├────────────────┬────────────────┬────────────────┐     │
     │ Airlines       │ Airports       │ Aircraft       │     │
     │ - airline_id   │ - airport_id   │ - aircraft_id  │     │
     │ - name         │ - code         │ - type_id      │     │
     │ - code         │ - name         │ - tail_number  │     │
     │ - country      │ - city         │ - airline_id   │     │
     │ - ...          │ - timezone     │ - ...          │     │
     │                │ - ...          │                │     │
     ├────────────────┼────────────────┼────────────────┤     │
     │ Flights        │ Schedules      │ Passengers     │     │
     │ - flight_id    │ - schedule_id  │ - passenger_id │     │
     │ - flight_number│ - flight_id    │ - name         │     │
     │ - airline_id   │ - departure    │ - email        │     │
     │ - route_id     │ - arrival      │ - ...          │     │
     │ - aircraft_id  │ - days         │                │     │
     │ - ...          │ - ...          │                │     │
     ├────────────────┴────────────────┴────────────────┤     │
     │ All services have access to complete database     │     │
     │ with permission controls on write operations      │     │
     └─────────────────────────────────────────────────────────┘
     ```
   - **Advantages**:
     - Easier initial migration with minimal schema changes
     - Familiar data model for developers
     - Simplified joins for complex queries
     - Reduced duplication of reference data
   - **Disadvantages**:
     - Tight coupling between services
     - Potential consistency issues when schema evolves
     - Requires coordination for schema changes
     - Database becomes a single point of failure

#### Long-term: Database per Service with Independent Schema
   - Each service exclusively owns its tables
   - Services communicate via well-defined APIs
   - Event sourcing for data synchronization
   - **Independent Schema Structure**:
     ```
     ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
     │ Airline Service │  │ Airport Service │  │ Fleet Service   │
     │ Database        │  │ Database        │  │ Database        │
     ├─────────────────┤  ├─────────────────┤  ├─────────────────┤
     │ - airlines      │  │ - airports      │  │ - aircraft      │
     │                 │  │ - timezones     │  │ - aircraft_types│
     │                 │  │                 │  │ - maintenance   │
     └─────────────────┘  └─────────────────┘  └─────────────────┘
                                         
     ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
     │ Flight Service  │  │ Booking Service │  │ Schedule Service│
     │ Database        │  │ Database        │  │ Database        │
     ├─────────────────┤  ├─────────────────┤  ├─────────────────┤
     │ - flights       │  │ - passengers    │  │ - schedules     │
     │ - routes        │  │ - bookings      │  │ - recurrences   │
     │ - status        │  │ - seat_maps     │  │                 │
     │ - ...           │  │ - ...           │  │ - ...           │
     └─────────────────┘  └─────────────────┘  └─────────────────┘
     ```
   - **Advantages**:
     - Loose coupling and independent deployment
     - Technology flexibility (different DB technologies per service)
     - Independent scaling based on service needs
     - Improved fault isolation
     - Schema changes don't affect other services
   - **Disadvantages**:
     - Increased complexity in data management
     - Eventual consistency challenges
     - Complex queries across service boundaries
     - Data duplication of reference/lookup data
     - More complex transaction management

## Communication Strategy

1. **Synchronous Communication**
   - RESTful APIs for direct client-server communication
   - gRPC for internal service-to-service communication (high performance)

2. **Asynchronous Communication**
   - Message broker (RabbitMQ/Kafka) for event-driven communication
   - Implement event sourcing for critical business events

## Simplified Migration Strategy

### Phase 1: Foundation (1 month)

1. **Simple API Gateway**
   - Create a basic API gateway to route requests between monolith and new services
   - Start with standard HTTP routing without complex orchestration

2. **Shared Library Extraction**
   - Extract common code into shared packages (authentication, validation, etc.)
   - Create simple database access patterns for the shared schema

### Phase 2: Initial Services (2 months)

1. **Extract Reference & Fleet Services**
   - Move Airport, Airline, and Fleet (aircraft) services first
   - These have fewer dependencies and are good candidates for early extraction
   - Implement with shared database access pattern

### Phase 3: Core Services (2-3 months)

**Note: These services can be migrated concurrently by separate teams**

1. **Extract Schedule & Flight Services**
   - Move schedule and flight management logic to dedicated services
   - Can be developed in parallel by separate teams

2. **Extract Booking Service**
   - Move passenger and booking management
   - Can be developed concurrently with schedule/flight services

3. **Add Notification Functionality**
   - Implement as part of booking service initially for simplicity
   - Later extract to dedicated service if volume justifies it

### Phase 4: Completion (1 month)

1. **Decommission Monolith**
   - Confirm all functionality migrated correctly
   - Remove redundant monolith code and route all traffic through services
   - Add service-specific scaling as needed

### Concurrent Development Approach

- Use feature branches for each service migration
- Each team can work independently on their service extraction
- Core services (flight, schedule, booking) can be developed in parallel
- Use the shared database schema to simplify initial development
- Regular integration tests to ensure services work together correctly

## Essential Testing & Operations

### Testing Approach
- Unit and integration tests for each service
- API contract testing to ensure service compatibility
- End-to-end testing for critical business flows (booking, scheduling)

### Simple Monitoring
- Basic logging across services with consistent format
- Health checks and basic metrics for each service
- Centralized error tracking (Sentry or similar)

### Security Essentials
- Service-to-service authentication with API keys or JWT
- Secure credentials storage
- Input validation and sanitization

### Scaling Approach
- Container-based scaling for high-traffic services
- Shared database with read replicas for performance
- Focus scaling efforts on booking/search during peak periods

## Timeline & Resources

- **Simplified Timeline**: 4-6 months total
- **Team Requirements**:
  - Backend Developers (3-4)
  - QA Engineer (1)
  - Project Manager (part-time)

## Success Metrics

- **Performance**: Response time & throughput improvements
- **Operational**: Ability to update/deploy services independently
- **Business**: Improved availability & capacity for peak periods