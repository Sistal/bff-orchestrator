# BFF Orchestrator

Backend for Frontend (BFF) service built with Go and Gin framework, implementing hexagonal architecture (Ports and Adapters pattern).

## Overview

This service acts as a BFF layer between a React frontend and multiple backend APIs. It orchestrates and aggregates data from various sources, providing a unified API tailored for frontend consumption.

## Architecture

The project follows **Hexagonal Architecture** (also known as Ports and Adapters), which provides:

- **Clear separation of concerns**: Business logic is isolated from external dependencies
- **Testability**: Easy to mock dependencies and test core logic
- **Flexibility**: Easy to swap implementations without affecting core business logic
- **Maintainability**: Well-organized code structure

### Project Structure

```
bff-orchestrator/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration management
│   ├── domain/
│   │   └── model/                     # Domain models
│   │       ├── user.go
│   │       ├── product.go
│   │       └── dashboard.go
│   ├── ports/                         # Interface definitions (ports)
│   │   ├── providers.go               # External service interfaces
│   │   └── services.go                # Business service interfaces
│   ├── application/
│   │   └── service/                   # Business logic (use cases)
│   │       └── aggregation_service.go
│   └── adapters/                      # Implementations (adapters)
│       ├── http/
│       │   ├── handler/               # HTTP handlers
│       │   ├── router/                # Route definitions
│       │   └── middleware/            # HTTP middleware (CORS, etc.)
│       └── provider/                  # External API clients
│           ├── user_provider.go
│           └── product_provider.go
├── Dockerfile                         # Container definition
├── .env.example                       # Environment variables template
├── go.mod                             # Go module definition
└── README.md                          # This file
```

### Hexagonal Architecture Layers

1. **Domain Layer** (`internal/domain`): Core business models and entities
2. **Ports Layer** (`internal/ports`): Interfaces that define contracts
3. **Application Layer** (`internal/application`): Business logic and use cases
4. **Adapters Layer** (`internal/adapters`): Concrete implementations
   - **HTTP Adapters**: REST API handlers and routes
   - **Provider Adapters**: External API clients

## Features

- ✅ RESTful API with Gin framework
- ✅ Hexagonal architecture for clean code organization
- ✅ CORS middleware for React frontend integration
- ✅ Health check endpoint
- ✅ Graceful shutdown
- ✅ Environment-based configuration
- ✅ Docker support
- ✅ No database dependency (stateless service)
- ✅ Example integrations with external APIs

## Prerequisites

- Go 1.23 or higher
- Docker (optional, for containerization)

## Getting Started

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/Sistal/bff-orchestrator.git
   cd bff-orchestrator
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the service**
   ```bash
   go run cmd/api/main.go
   ```

   The service will start on `http://localhost:8080`

### Using Docker

1. **Build the Docker image**
   ```bash
   docker build -t bff-orchestrator:latest .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 \
     -e PORT=8080 \
     -e ENVIRONMENT=production \
     -e USER_SERVICE_URL=https://jsonplaceholder.typicode.com \
     -e PRODUCT_SERVICE_URL=https://fakestoreapi.com \
     bff-orchestrator:latest
   ```

## API Endpoints

### Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "bff-orchestrator"
}
```

### Dashboard (Aggregated Data)
```
GET /api/v1/dashboard/:userId
```

**Description:** Aggregates user and product data from multiple external APIs.

**Example Request:**
```bash
curl http://localhost:8080/api/v1/dashboard/1
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "name": "Leanne Graham",
    "email": "Sincere@april.biz",
    "username": "Bret"
  },
  "products": [
    {
      "id": 1,
      "title": "Product Name",
      "price": 29.99,
      "description": "Product description",
      "category": "electronics"
    }
  ]
}
```

## Configuration

Configure the service using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `USER_SERVICE_URL` | External user service URL | `https://jsonplaceholder.typicode.com` |
| `PRODUCT_SERVICE_URL` | External product service URL | `https://fakestoreapi.com` |

## Adding New Features

### Adding a New External Service

1. **Define the port** in `internal/ports/providers.go`:
   ```go
   type OrderProvider interface {
       GetOrders(ctx context.Context, userID int) ([]*model.Order, error)
   }
   ```

2. **Create the adapter** in `internal/adapters/provider/`:
   ```go
   type OrderProvider struct {
       baseURL    string
       httpClient *http.Client
   }
   ```

3. **Inject into services** via constructor in `cmd/api/main.go`

### Adding a New Endpoint

1. **Create handler** in `internal/adapters/http/handler/`
2. **Register route** in `internal/adapters/http/router/router.go`
3. **Wire dependencies** in `cmd/api/main.go`

## Development

### Build
```bash
go build -o bin/bff-orchestrator ./cmd/api
```

### Run tests (when added)
```bash
go test ./...
```

### Format code
```bash
go fmt ./...
```

### Lint code (requires golangci-lint)
```bash
golangci-lint run
```

## Production Deployment

The service is designed to be deployed as a stateless container:

1. Build the Docker image
2. Deploy to your container orchestration platform (Kubernetes, ECS, etc.)
3. Configure environment variables
4. Set up health check monitoring using `/health` endpoint
5. Configure load balancing and auto-scaling as needed

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Ensure code quality and add tests
5. Submit a pull request

## 📚 API Contracts Documentation

Este BFF consume datos de **4 microservicios** backend. La documentación completa de los contratos de API se encuentra en la carpeta [`api-contracts/`](api-contracts/).

### Documentación Disponible

| Documento | Descripción |
|-----------|-------------|
| [**INDEX.md**](api-contracts/INDEX.md) | 📖 Índice principal y guía general |
| [**MS-Authentication**](api-contracts/MS-AUTHENTICATION-CONTRACT.md) | 🔐 Servicio de autenticación JWT (Puerto 8081) |
| [**MS-Funcionario**](api-contracts/MS-FUNCIONARIO-CONTRACT.md) | 👥 Servicio de recursos humanos (Puerto 8082) |
| [**MS-Operations**](api-contracts/MS-OPERATIONS-CONTRACT.md) | 📦 Servicio de operaciones (Puerto 8083) |
| [**MS-Catalog**](api-contracts/MS-CATALOG-CONTRACT.md) | 📋 Servicio de catálogos maestros (Puerto 8084) |
| [**DATABASE-SCHEMA.md**](api-contracts/DATABASE-SCHEMA.md) | 🗄️ Esquema completo de base de datos |
| [**IMPLEMENTATION-GUIDE.md**](api-contracts/IMPLEMENTATION-GUIDE.md) | 🛠️ Guía de implementación |
| [**USAGE-EXAMPLES.md**](api-contracts/USAGE-EXAMPLES.md) | 💻 Ejemplos de uso (cURL, JS, Go) |
| [**DIAGRAMS.md**](api-contracts/DIAGRAMS.md) | 📊 Diagramas de flujo y arquitectura |

### Arquitectura de Microservicios

```
┌─────────────────┐         ┌─────────────────┐
│   Frontend 1    │         │   Frontend 2    │
│  (Employee App) │         │  (Admin Panel)  │
└────────┬────────┘         └────────┬────────┘
         │                           │
         └───────────┬───────────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │   BFF Orchestrator    │
         │    (Port 8080)        │
         └───────────┬───────────┘
                     │
         ┌───────────┼───────────┬───────────┐
         │           │           │           │
         ▼           ▼           ▼           ▼
┌────────────┐ ┌──────────┐ ┌──────────┐ ┌─────────┐
│MS-Auth     │ │MS-Funcio │ │MS-Opera  │ │MS-Catalog│
│(Port 8081) │ │(Port 8082)│ │(Port 8083)│ │(Port 8084)│
└──────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬────┘
       │            │            │            │
       └────────────┴────────────┴────────────┘
                    │
                    ▼
         ┌──────────────────────┐
         │  PostgreSQL (Supabase)│
         └──────────────────────┘
```

### Inicio Rápido

1. **Revisar contratos**: Comienza con el [INDEX.md](api-contracts/INDEX.md)
2. **Configurar BD**: Usa el [DATABASE-SCHEMA.md](api-contracts/DATABASE-SCHEMA.md)
3. **Implementar servicios**: Sigue la [IMPLEMENTATION-GUIDE.md](api-contracts/IMPLEMENTATION-GUIDE.md)
4. **Probar endpoints**: Revisa los [USAGE-EXAMPLES.md](api-contracts/USAGE-EXAMPLES.md)

---

## License

MIT License

## Contact

For questions or support, please open an issue in the repository.