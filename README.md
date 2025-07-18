# Fast Food FIAP Tech Challenge 4 - Customer Service

## 💬 About

This project implements a serverless authentication service using Go, Clean Architecture, AWS Lambda and AWS API
Gateway. The service receives customer credentials, validates them, and returns a signed JWT token upon successful
authentication. The architecture enables scalability, maintainability, and testability.

## 🔗 Related Projects

This project is part of a larger system that includes:

- [Database Infrastructure (Terraform)](https://github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-db-tf) - Infrastructure
  as Code for MongoDB Atlas using Terraform
- [Kubernetes Infrastructure (Terraform)](https://github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-k8s-tf) -
  Infrastructure as Code for EKS cluster and Kubernetes resources using Terraform
- [API Service](https://github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-api) - Main backend service implementing the Fast
  Food ordering system

---

## 📁 Folder Structure

```bash
├── bootstrap
├── docs
│ └── architecture.drawio
├── internal
│ ├── adapter
│ │ ├── controller
│ │ ├── gateway
│ │ └── presenter
│ ├── core
│ │ ├── domain
│ │ │ ├── entity
│ │ │ └── errors.go
│ │ ├── dto
│ │ ├── port
│ │ │ └── mocks
│ │ └── usecase
│ └── infrastructure
│     ├── aws
│     │ └── lambda
│     │     ├── golden
│     │     ├── request
│     │     └── response
│     ├── config
│     ├── database
│     ├── datasource
│     ├── logger
│     └── service
├── terraform
│ ├── modules
│ │   ├── apigateway
│ │   └── lambda
│   └── test
└── fixture
```

---

## 🚀 Features

- Customer authentication via email and password
- Secure JWT generation for authenticated sessions
- Clean Architecture separation (domain, use cases, adapters, infrastructure)
- Unit tests with testify and golden file responses
- Error response standardization
- Environment-based configuration
- Terraform for AWS Lambda, API Gateway, IAM provisioning

---

## 🔧 Technologies

- **Go**
- **AWS Lambda**
- **Terraform**
- **Docker**
- **Docker Compose**
- **MongoDB**
- **Testify**
- **JWT**
- **Makefile**
- **Structured logging**

---

## ⚙️ Getting Started

### Prerequisites

- Go 1.24+
- AWS CLI
- Terraform
- Docker
- MongoDB (for local development)

### Local Development

1. Clone the repository:

   ```bash
      git clone https://github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf.git
      cd fiap-tech-challenge-3-lambda-auth-tf
   ```

2. Create your environment variables:

   ```shell
   cp env.example .env
   # Edit .env as needed 
   ```

3. Install dependencies:

   ```shell
   make install
   ```

4. Initialize lambda to receive requests:

   ```shell
   # Starts database
   make compose-up
   # Starts lambda
   make start-lambda
   ```

5. Trigger lambda events

   ```shell
   make trigger-lambda 
   ```

6. Run tests

   ```shell
   make test 
   ```

7. View coverage:

   ```shell
   make coverage
   ```

## 📝 Authentication API

## 🏗️ Deployment

Deployment is automated via a **GitHub Actions workflow**. When changes are pushed to the main branch (or as configured
in your workflow), the pipeline will build and deploy the Lambda function and related infrastructure using Terraform.

**Prerequisite:**
Before running `terraform plan` or `terraform apply` (either locally or via CI), ensure that all variables defined in
`terraform/modules/lambda/ssm.tf` are created and initialized in your AWS environment. These variables are required for
successful provisioning and configuration of the Lambda function and related resources.

All the variables can be found on `env.example` file.

## 📈 Testing

Unit tests: make test
Coverage: make coverage
Golden files for output validation are found in internal/infrastructure/aws/lambda/golden/.

## 🧩 Architecture

The project follows Clean Architecture, dividing source code into distinct layers: Domain, UseCases, Adapters, and
Infrastructure. See docs/architecture.drawio for the infrastructure diagram.

## 👏 Contributing

Fork the repository and create your branch from master branch.
Run tests before PR (make test)
Ensure code style with make lint
Follow Conventional Commits for commit messages

## 🙏 Support

For issues, open a GitHub issue in this repository.

## 📚 Docs

- [Best practices writing lambda functions](https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html)
- [Code best practices for Go Lambda functions](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#go-best-practices)
- [Running and debugging lambda locally](https://medium.com/nagoya-foundation/running-and-debugging-go-lambda-functions-locally-156893e4ed0d)
- [Setting Up VPC and Lambda Function with Terraform](https://dev.to/sepiyush/setting-up-vpc-and-lambda-function-with-terraform-3m9d)
- [MongoDB Go Driver Documentation](https://www.mongodb.com/docs/drivers/go/current/)
- [MongoDB Best Practices](https://www.mongodb.com/developer/products/mongodb/mongodb-schema-design-best-practices/)

## 📄 License

MIT License
