# Library Service 📚 [![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](https://go.dev/) [![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=fff)](https://www.docker.com/) [![Postgres](https://img.shields.io/badge/Postgres-%23316192.svg?logo=postgresql&logoColor=white)](https://www.postgresql.org/)

## 🚀 Project Overview

This project demonstrates the difference between directly dependent and dependency-injected architectures in Go backend applications. The repository contains two different implementations of the same library management service:

1. **Direct Service** - Uses direct dependencies with package-level globals and function calls
2. **Injected Service** - Uses dependency injection pattern with interfaces

## 🏗️ Architecture Comparison

### Direct Service

The direct service implementation:
- Direct function calls
- Tight coupling between components
- Hard to replace or mock dependencies
- Currently only tested via integration test

### Injected Service

The injected service implementation:
- Dependency injection via interfaces
- Loose coupling between components
- Easy testing through dependency mocking
- Simple component replacement
- Near 100% test coverage

## ✨ Key Features

- Complete library management system with books, users, and lending functionality
- RESTful API implementation for all CRUD operations
- PostgreSQL database integration
- Docker containerization for easy deployment

## 🏁 Getting Started

### Prerequisites

- ![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
- ![Docker](https://img.shields.io/badge/Docker-Required-2496ED?style=flat&logo=docker)
- ![Mockery](https://img.shields.io/badge/Mockery-Optional-orange?style=flat&logo=mockery)

### Running the Application

1. Clone the repository
2. Start the services using Docker Compose:

```bash
docker-compose up
```

This will start:
- PostgreSQL database on port 5432
- Direct service on port 8081
- Injected service on port 8080
- Database migration service

## 🧪 Testing and Coverage

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