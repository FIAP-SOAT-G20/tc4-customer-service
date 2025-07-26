# Fast Food FIAP Tech Challenge 4 - Customer Service

[![Tests](https://github.com/FIAP-SOAT-G20/tc4-customer-service/workflows/Tests/badge.svg)](https://github.com/FIAP-SOAT-G20/tc4-customer-service/actions/workflows/test.yml)
[![Build and Deploy](https://github.com/FIAP-SOAT-G20/tc4-customer-service/workflows/Build%20and%20Deploy/badge.svg)](https://github.com/FIAP-SOAT-G20/tc4-customer-service/actions/workflows/build-deploy.yml)
[![Go Version](https://img.shields.io/badge/go-1.24.2-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 💬 Visão Geral

Este projeto implementa um serviço de autenticação e gerenciamento de clientes serverless usando Go, Clean Architecture,
AWS Lambda e AWS API Gateway. O serviço recebe credenciais de clientes, valida-as e retorna um token JWT assinado após
autenticação bem-sucedida. A arquitetura permite escalabilidade, manutenibilidade e testabilidade.

### Principais Funcionalidades

- Autenticação de clientes via email e senha
- Geração segura de JWT para sessões autenticadas
- CRUD completo de clientes
- Separação por Clean Architecture (domain, use cases, adapters, infrastructure)
- Testes unitários com testify e golden file responses
- Padronização de respostas de erro
- Configuração baseada em ambiente

---

## 🏗️ Tecnologias e Estrutura

### Estrutura do Projeto

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
├── terraform/              # Infrastructure as Code
│ └── modules/              # Terraform modules
└── test/                   # Test data and fixtures
```

### Tecnologias

- **Go 1.24.2** - Linguagem de programação
- **AWS Lambda** - Plataforma serverless
- **AWS API Gateway** - Gerenciamento de API
- **Amazon ECR** - Registry de containers
- **MongoDB** - Banco de dados NoSQL
- **Terraform** - Infrastructure as Code
- **Docker** - Containerização
- **GitHub Actions** - Pipeline CI/CD
- **JWT** - Tokens de autenticação
- **Testify** - Framework de testes
- **golangci-lint** - Linting de código

---

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.24.2+
- AWS CLI
- Terraform
- Docker
- MongoDB

### Instalação e Execução

1. **Clone o repositório:**

   ```bash
   git clone https://github.com/FIAP-SOAT-G20/tc4-customer-service.git
   cd tc4-customer-service
   ```

2. **Configure as variáveis de ambiente:**

   ```bash
   cp env.example .env
   # Edite .env conforme necessário 
   ```

3. **Instale as dependências:**

   ```bash
   make install
   ```

4. **Inicie o ambiente de desenvolvimento:**

   ```bash
   # Inicia o banco de dados
   make compose-up
   
   # Inicia o lambda
   make start-lambda
   ```

5. **Teste o lambda:**

   ```bash
   make trigger-lambda 
   ```

### Comandos Disponíveis

```bash
make help          # Mostra todos os comandos disponíveis
make build         # Compila a aplicação
make test          # Executa os testes
make coverage      # Gera relatório de coverage
make lint          # Executa linter
make scan          # Executa scan de segurança
make package       # Empacota para deploy
make compose-up    # Inicia ambiente local
make compose-down  # Para ambiente local
```

---

## 📝 API e Documentação

### Endpoints Disponíveis

| Método   | Endpoint               | Descrição                           |
|----------|------------------------|-------------------------------------|
| `POST`   | `/auth`                | Autentica cliente com email e senha |
| `GET`    | `/customers/{id}`      | Busca cliente por ID                |
| `GET`    | `/customers/cpf/{cpf}` | Busca cliente por CPF               |
| `GET`    | `/customers`           | Lista todos os clientes             |
| `POST`   | `/customers`           | Cria novo cliente                   |
| `PUT`    | `/customers/{id}`      | Atualiza cliente                    |
| `DELETE` | `/customers/{id}`      | Remove cliente                      |

---

## 🧪 Testes e Qualidade

### Execução de Testes

```bash
# Executa todos os testes com detecção de race condition
make test

# Gera relatório de coverage (abre no browser)
make coverage

# Executa linter
make lint

# Executa scan de vulnerabilidades
make scan
```

## 🏗️ Deploy e CI/CD

### Pipeline Automatizado

1. **Tests workflow** - Executa em todo push/PR para branch main
2. **Build and Deploy workflow** - Dispara após testes bem-sucedidos

### Processo de Deploy

- **Testes**: Testes automatizados com upload de coverage para Codecov
- **Linting**: Verificações de qualidade de código
- **Security**: Scan de vulnerabilidades
- **Build**: Criação de imagem Docker e push para ECR
- **Deploy**: Deploy automatizado para AWS Lambda

### Pré-requisitos para Deploy

- Credenciais AWS configuradas nos GitHub Secrets
- Repositório ECR será criado automaticamente se não existir

### Deploy Local com Terraform

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

---

## 🔗 Projetos Relacionados

Este projeto faz parte de um sistema maior que inclui:

- **[Infrastructure (Terraform)](https://github.com/FIAP-SOAT-G20/tc4-infrastructure-tf)** - Infrastructure as Code para
  recursos AWS
- **[Customer Service](https://github.com/FIAP-SOAT-G20/tc4-customer-service)** - Serviço de autenticação e
  gerenciamento de clientes
- **[Payment Service](https://github.com/FIAP-SOAT-G20/tc4-payment-service)** - Serviço de processamento de pagamentos
- **[Kitchen Service](https://github.com/FIAP-SOAT-G20/tc4-kitchen-service)** - Serviço de operações da cozinha e
  gerenciamento de pedidos
- **[Kubernetes Deploy](https://github.com/FIAP-SOAT-G20/tc4-infrastructure-deploy)** - Configurações de deploy no
  Kubernetes

---

## 📚 Documentação de Referência

- [Best practices writing lambda functions](https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html)
- [Code best practices for Go Lambda functions](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#go-best-practices)
- [Running and debugging lambda locally](https://medium.com/nagoya-foundation/running-and-debugging-go-lambda-functions-locally-156893e4ed0d)
- [Setting Up VPC and Lambda Function with Terraform](https://dev.to/sepiyush/setting-up-vpc-and-lambda-function-with-terraform-3m9d)
- [MongoDB Go Driver Documentation](https://www.mongodb.com/docs/drivers/go/current/)
- [MongoDB Best Practices](https://www.mongodb.com/developer/products/mongodb/mongodb-schema-design-best-practices/)

## 📄 Licença

MIT License