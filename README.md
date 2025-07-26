# Fast Food FIAP Tech Challenge 4 - Customer Service

[![Tests](https://github.com/FIAP-SOAT-G20/tc4-customer-service/workflows/Tests/badge.svg)](https://github.com/FIAP-SOAT-G20/tc4-customer-service/actions/workflows/test.yml)
[![Build and Deploy](https://github.com/FIAP-SOAT-G20/tc4-customer-service/workflows/Build%20and%20Deploy/badge.svg)](https://github.com/FIAP-SOAT-G20/tc4-customer-service/actions/workflows/build-deploy.yml)
[![Go Version](https://img.shields.io/badge/go-1.24.2-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 💬 Overview

This project implements a serverless customer authentication and management service using Go, Clean Architecture,
and AWS Lambda. The service receives customer credentials, validates them, and returns a signed JWT token after
successful authentication. The architecture enables scalability, maintainability, and testability.

### Key Features

- Customer authentication via email and password
- Secure JWT generation for authenticated sessions
- Complete customer CRUD operations
- Clean Architecture separation (domain, use cases, adapters, infrastructure)
- Unit tests with testify and golden file responses
- Standardized error responses
- Environment-based configuration

---

## 🏗️ Technologies and Structure

### Project Structure

```bash
├── bin/                    # Compiled binaries
├── dist/                   # Distribution files
├── internal/               # Private application code
│ ├── adapter/              # External interface adapters
│ │ ├── controller/         # HTTP handlers
│ │ ├── gateway/            # External service interfaces
│ │ └── presenter/          # Response formatters
│ ├── core/                 # Business logic
│ │ ├── domain/             # Domain entities and rules
│ │ ├── dto/                # Data transfer objects
│ │ ├── port/               # Interfaces and mocks
│ │ └── usecase/            # Business use cases
│ └── infrastructure/       # External concerns
│     ├── aws/lambda/       # AWS Lambda integration
│     ├── config/           # Configuration management
│     ├── database/         # Database connections
│     ├── datasource/       # Data access layer
│     ├── logger/           # Logging utilities
│     └── service/          # External services
└── test/                   # Test data and fixtures
```

### Technologies

- **Go 1.24.2** - Programming language
- **AWS Lambda** - Serverless platform
- **Amazon ECR** - Container registry
- **MongoDB** - NoSQL database
- **Docker** - Containerization
- **GitHub Actions** - CI/CD pipeline
- **JWT** - Authentication tokens
- **Testify** - Testing framework
- **golangci-lint** - Code linting

---

## 🚀 Quick Start

### Prerequisites

- Go 1.24.2+
- AWS CLI
- Docker
- MongoDB

### Installation and Execution

1. **Clone the repository:**

   ```bash
   git clone https://github.com/FIAP-SOAT-G20/tc4-customer-service.git
   cd tc4-customer-service
   ```

2. **Configure environment variables:**

   ```bash
   cp env.example .env
   # Edit .env as needed
   ```

3. **Install dependencies:**

   ```bash
   make install
   ```

4. **Start development environment:**

   ```bash
   # Start database
   make compose-up
   
   # Start lambda
   make start-lambda
   ```

5. **Test the lambda:**

   ```bash
   make trigger-lambda 
   ```

### Lambda Trigger Commands

You can trigger different lambda endpoints using predefined test events:

```bash
# Default trigger (customer not found scenario)
make trigger-lambda

# Authentication
LAMBDA_INPUT_FILE=test/data/auth_customer.json make trigger-lambda

# Customer CRUD operations
LAMBDA_INPUT_FILE=test/data/create_customer.json make trigger-lambda
LAMBDA_INPUT_FILE=test/data/get_customer_by_id.json make trigger-lambda
LAMBDA_INPUT_FILE=test/data/get_customer_by_cpf.json make trigger-lambda
LAMBDA_INPUT_FILE=test/data/update_customer.json make trigger-lambda
LAMBDA_INPUT_FILE=test/data/delete_customer.json make trigger-lambda
LAMBDA_INPUT_FILE=test/data/list_customers.json make trigger-lambda

# Edge cases
LAMBDA_INPUT_FILE=test/data/api_gateway_proxy_request_event_payload_empty_cpf.json make trigger-lambda
```

### Available Commands

```bash
make help          # Show all available commands
make build         # Build the application
make test          # Run tests
make coverage      # Generate coverage report
make lint          # Run linter
make scan          # Run security scan
make package       # Package for deployment
make compose-up    # Start local environment
make compose-down  # Stop local environment
```

---

## 📝 API Documentation

### Available Endpoints

| Method   | Endpoint               | Description                                   |
|----------|------------------------|-----------------------------------------------|
| `POST`   | `/auth`                | Authenticate customer with email and password |
| `GET`    | `/customers/{id}`      | Get customer by ID                            |
| `GET`    | `/customers/cpf/{cpf}` | Get customer by CPF                           |
| `GET`    | `/customers`           | List all customers                            |
| `POST`   | `/customers`           | Create new customer                           |
| `PUT`    | `/customers/{id}`      | Update customer                               |
| `DELETE` | `/customers/{id}`      | Delete customer                               |

---

## 🧪 Testing and Quality

### Running Tests

```bash
# Run all tests with race condition detection
make test

# Generate coverage report (opens in browser)
make coverage

# Run linter
make lint

# Run vulnerability scan
make scan
```

## 🏗️ Deploy and CI/CD

### Automated Pipeline

1. **Tests workflow** - Runs on every push/PR to main branch
2. **Build and Deploy workflow** - Triggers after successful tests

### Deploy Process

- **Testing**: Automated tests with coverage upload to Codecov
- **Linting**: Code quality checks
- **Security**: Vulnerability scanning
- **Build**: Docker image creation and push to ECR
- **Deploy**: Automated deployment to AWS Lambda

### Deploy Prerequisites

- AWS credentials configured in GitHub Secrets
- ECR repository will be created automatically if it doesn't exist

---

## 🔗 Related Projects

This project is part of a larger system that includes:

- **[Infrastructure (Terraform)](https://github.com/FIAP-SOAT-G20/tc4-infrastructure-tf)** - Infrastructure as Code for
  AWS resources
- **[Customer Service](https://github.com/FIAP-SOAT-G20/tc4-customer-service)** - Customer authentication and
  management service
- **[Payment Service](https://github.com/FIAP-SOAT-G20/tc4-payment-service)** - Payment processing service
- **[Kitchen Service](https://github.com/FIAP-SOAT-G20/tc4-kitchen-service)** - Kitchen operations and
  order management service
- **[Kubernetes Deploy](https://github.com/FIAP-SOAT-G20/tc4-infrastructure-deploy)** - Kubernetes deployment
  configurations

---

## 📚 Reference Documentation

- [Best practices writing lambda functions](https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html)
- [Code best practices for Go Lambda functions](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#go-best-practices)
- [Running and debugging lambda locally](https://medium.com/nagoya-foundation/running-and-debugging-go-lambda-functions-locally-156893e4ed0d)
- [MongoDB Go Driver Documentation](https://www.mongodb.com/docs/drivers/go/current/)
- [MongoDB Best Practices](https://www.mongodb.com/developer/products/mongodb/mongodb-schema-design-best-practices/)

## 📄 License

MIT License