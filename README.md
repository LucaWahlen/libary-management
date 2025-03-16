# Library Service

## Project Overview

This project demonstrates the difference between directly dependent and dependency-injected architectures in Go backend applications. The repository contains two different implementations of the same library management service:

1. **Direct Service** - Uses direct dependencies with package-level globals and function calls
2. **Injected Service** - Uses dependency injection pattern with interfaces

By comparing these two approaches side by side, this project illustrates how dependency injection impacts code quality, testability, maintainability, and scalability in large-scale Go applications.

## Architecture Comparison

### Direct Service

The direct service implementation:
- Uses global package variables and direct function calls
- Has tight coupling between components
- Makes testing more difficult due to hidden dependencies
- Requires more effort to replace or mock dependencies
- Currently only tested with the integration test, but without coverage

```go
// Example from direct-service
func CreateBook(w http.ResponseWriter, r *http.Request) {
    // Direct dependency on repository package
    createdBook, err := repository.CreateBook(book)
    // ...
}
```

### Injected Service

The injected service implementation:
- Uses dependency injection via interfaces
- Has loose coupling between components
- Facilitates easier testing through dependency mocking
- Allows for straightforward component replacement
- Currently near 100% test coverage with integration & unit tests

```go
// Example from injected-service
func (s *LibaryService) CreateBook(w http.ResponseWriter, r *http.Request) {
    // Injected dependency used via interface
    createdBook, err := s.repository.CreateBook(book)
    // ...
}
```

## Project Structure

```
├── cmd
│   ├── direct-service        # Direct dependency implementation
│   └── injected-service      # Dependency injection implementation
├── internal
│   ├── direct-service        # Implementation with direct dependencies
│   │   ├── app               # Application logic with direct dependencies
│   │   ├── repository        # Data access with global variables
│   │   ├── router            # HTTP routing with direct function calls
│   │   └── validation        # Validation logic with direct dependencies
│   ├── domain                # Shared domain models
│   └── injected-service      # Implementation with dependency injection
│       ├── app               # Application logic with injected dependencies
│       ├── repository        # Repository interfaces and implementations
│       ├── router            # HTTP routing with injected dependencies
│       └── validation        # Validation with injected dependencies
├── integrationtest           # Integration tests for both services
├── migrations                # Database migration files
└── docker-compose.yaml       # Docker configuration for running services
```

## Key Features

- Complete library management system with books, users, and lending functionality
- RESTful API implementation for all CRUD operations
- PostgreSQL database integration
- Docker containerization for easy deployment
- Integration tests that validate both implementations

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker and Docker Compose
- PostgreSQL client (for direct DB interaction, optional)
- Mockery

### Running the Application

1. Clone the repository
2. Start the services using Docker Compose:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Direct service on port 8081
- Injected service on port 8080
- Database migration service

```bash
./runlocal-testcoverage.sh
```


## Testing and Coverage

Before running tests, make sure to generate the necessary mock implementations by executing:

```bash
go generate ./...
```
Run tests with coverage report:

```bash
./runlocal-testcoverage.sh
```

This script will:
1. Run all tests
2. Generate a coverage report
3. Open the report in your default browser