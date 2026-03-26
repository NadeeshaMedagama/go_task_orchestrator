# Go Task Orchestrator

A production-ready microservices task management platform built with **Go** backend services and a **Next.js** frontend. The project showcases modern cloud-native architecture with containerized services, JWT authentication, and a comprehensive CI/CD pipeline.

## 🏗️ Architecture Overview

### Services
- **API Gateway** (`cmd/api-gateway:8080`) - Single entry point, JWT validation, request routing
- **Auth Service** (`cmd/auth-service:8081`) - User registration/login, JWT token issuance, user management (RBAC)
- **Task Service** (`cmd/task-service:8082`) - Task CRUD, status management, filtering, pagination, RBAC
- **Frontend** (`frontend:3000`) - Next.js React app with modern UI/UX

### Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Backend | Go | 1.21+ |
| API Pattern | REST with JWT | - |
| Frontend Framework | Next.js | 16.2.1 |
| Frontend UI | React | 19+ |
| Styling | Tailwind CSS | 4.2.1 |
| HTTP Client | Axios | 1.6.7+ |
| Container Runtime | Docker & Docker Compose | Latest |
| CI/CD | GitHub Actions | - |
| Container Orchestration | Kubernetes | 1.25+ (optional) |

## 📁 Repository Structure

```
Go_Task_Orchestrator/
├── cmd/
│   ├── api-gateway/          # Main API gateway service
│   ├── auth-service/         # Authentication service
│   └── task-service/         # Task management service
├── frontend/                 # Next.js React frontend
│   ├── src/
│   │   ├── app/              # Next.js app directory
│   │   ├── components/       # Reusable React components
│   │   ├── hooks/            # Custom React hooks
│   │   ├── lib/              # Utilities (auth, axios config)
│   │   ├── context/          # React context (authentication)
│   │   └── types/            # TypeScript definitions
│   ├── Dockerfile            # Multi-stage Docker build
│   └── package.json          # npm dependencies
├── internal/
│   └── platform/
│       └── jwtutil/          # Shared JWT utilities
├── deployments/
│   └── docker/               # Dockerfiles for each service
├── docs/                     # Project documentation
├── k8s/                      # Kubernetes manifests
├── .github/
│   └── workflows/            # CI/CD pipelines
├── docker-compose.yml        # Local development orchestration
├── go.mod & go.sum           # Go module dependencies
└── action.yml                # GitHub Action definition

```

## 🚀 Quick Start

### Prerequisites
- **Docker & Docker Compose** — For containerized deployment
- **Go 1.21+** — For local backend development
- **Node.js 20+** — For local frontend development
- **Git** — For version control

### Option 1: Docker Compose (Recommended)

Launch the entire stack with a single command:

```bash
# Clone the repository
git clone <repo-url>
cd Go_Task_Orchestrator

# Copy environment configuration
cp .env.example .env

# Start all services
docker compose up --build
```

**Access the application:**
- **Frontend:** http://localhost:3000
- **API Gateway:** http://localhost:8080
- **Health Check:** http://localhost:8080/health

**Default Admin Credentials:**
- Username: `admin`
- Password: `admin123`

### Option 2: Local Development (Without Docker)

#### Start Backend Services

Open three separate terminal windows and run:

```bash
# Terminal 1: Auth Service
go run ./cmd/auth-service

# Terminal 2: Task Service
go run ./cmd/task-service

# Terminal 3: API Gateway
go run ./cmd/api-gateway
```

#### Start Frontend

```bash
cd frontend
npm install
npm run dev
```

**Access:**
- Frontend: http://localhost:3000
- API Gateway: http://localhost:8080

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```bash
POST /api/auth/register
Content-Type: application/json

{
  "username": "john",
  "email": "john@example.com",
  "password": "secure_password_123"
}
```

**Response:**
```json
{
  "id": "user-uuid",
  "username": "john",
  "email": "john@example.com",
  "token": "eyJhbGc...",
  "created_at": "2024-03-21T10:00:00Z"
}
```

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "john",
  "password": "secure_password_123"
}
```

**Response:**
```json
{
  "id": "user-uuid",
  "username": "john",
  "token": "eyJhbGc...",
  "created_at": "2024-03-21T10:00:00Z"
}
```

### Task Endpoints

#### Create Task
```bash
POST /api/tasks
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "title": "Implement feature X",
  "description": "Add new authentication flow",
  "priority": "HIGH",
  "status": "TODO"
}
```

#### Get All Tasks
```bash
GET /api/tasks?page=1&limit=10&status=TODO
Authorization: Bearer <TOKEN>
```

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)
- `status` - Filter by status (TODO, IN_PROGRESS, DONE, CANCELLED)
- `priority` - Filter by priority (LOW, MEDIUM, HIGH, URGENT)

#### Get Task by ID
```bash
GET /api/tasks/{id}
Authorization: Bearer <TOKEN>
```

#### Update Task
```bash
PUT /api/tasks/{id}
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "title": "Updated title",
  "status": "IN_PROGRESS",
  "priority": "MEDIUM"
}
```

#### Delete Task
```bash
DELETE /api/tasks/{id}
Authorization: Bearer <TOKEN>
```

#### List Users (Admin Only)
```bash
GET /api/users
Authorization: Bearer <ADMIN_TOKEN>
```

## 🔒 Frontend & Backend Connection

The frontend communicates with the backend through the API Gateway using **HTTP/REST with Bearer token authentication**.

### Configuration Files

**Frontend API Configuration** (`frontend/src/lib/axios.ts`):
- Base URL: `${NEXT_PUBLIC_API_URL}` (defaults to `http://localhost:8080/api`)
- Request interceptor automatically adds Bearer token from localStorage
- Response interceptor handles token expiration (clears auth on 401)

**Environment Variables:**

```env
# .env (local development)
NEXT_PUBLIC_API_URL=http://localhost:8080/api

# docker-compose.yml (Docker)
NEXT_PUBLIC_API_URL=http://api-gateway:8080/api
```

### Authentication Flow

1. User submits login/register form
2. Frontend sends credentials to `/api/auth/login` or `/api/auth/register`
3. Backend returns JWT token
4. Frontend stores token in localStorage
5. All subsequent requests include `Authorization: Bearer <token>` header
6. API Gateway validates JWT before routing to services
7. If token is invalid/expired, 401 response clears localStorage and redirects to login

### Tested Endpoints
✅ **Register** - Working
✅ **Login** - Working  
✅ **Get Tasks** - Working (with authentication)
✅ **Create Task** - Working
✅ **Update Task** - Working
✅ **Delete Task** - Working

## 🔧 Configuration

### Environment Variables

**Backend Services:**
```env
JWT_SECRET=<base64-encoded-jwt-secret>  # Required
PORT=8080                                # Default: 8080 (gateway), 8081 (auth), 8082 (task)
AUTH_SERVICE_URL=http://localhost:8081   # Gateway only
TASK_SERVICE_URL=http://localhost:8082   # Gateway only
```

**Frontend:**
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api  # Browser-accessible API
NEXT_TELEMETRY_DISABLED=1                      # Disable telemetry
```

**Generate JWT Secret:**
```bash
openssl rand -base64 32
```

## 🐳 Docker & Kubernetes

### Docker Compose Services

| Service | Container | Port | Dockerfile |
|---------|-----------|------|------------|
| api-gateway | go_task_orchestrator_gateway | 8080 | deployments/docker/api-gateway.Dockerfile |
| auth-service | go_task_orchestrator_auth | 8081 | deployments/docker/auth-service.Dockerfile |
| task-service | go_task_orchestrator_task | 8082 | deployments/docker/task-service.Dockerfile |
| frontend | go_task_orchestrator_frontend | 3000 | frontend/Dockerfile |

All services run on a shared `orchestrator-net` bridge network.

### Run with Kubernetes (Minikube)

To run this project locally on Kubernetes with Minikube, use the steps below.

```bash
# 1) Start Minikube
minikube start

# 2) Enable ingress
minikube addons enable ingress

# 3) Build images inside Minikube's Docker daemon
eval "$(minikube docker-env)"
docker build -t go-task-orchestrator/auth-service:local -f deployments/docker/auth-service.Dockerfile .
docker build -t go-task-orchestrator/task-service:local -f deployments/docker/task-service.Dockerfile .
docker build -t go-task-orchestrator/api-gateway:local -f deployments/docker/api-gateway.Dockerfile .

# 4) Deploy manifests with Kustomize
kubectl apply -k k8s

# 5) Watch pods until they are Running
kubectl get pods -n go-task-orchestrator -w
```

The Kubernetes manifests are configured with `imagePullPolicy: Never`, so local images built inside Minikube are used directly.

**Expose ingress host locally:**

```bash
minikube ip
```

Add this entry to `/etc/hosts` (replace with your Minikube IP):

```text
<MINIKUBE_IP> go-task-orchestrator.local
```

Then verify:

- Health: `http://go-task-orchestrator.local/health`
- API base: `http://go-task-orchestrator.local/api`

#### Run Frontend Against Minikube Backend

Use this when backend services are running in Minikube and you want to run the Next.js frontend locally.

```bash
# Ensure the ingress host resolves locally
MINIKUBE_IP=$(minikube ip)
echo "$MINIKUBE_IP go-task-orchestrator.local" | sudo tee -a /etc/hosts

# Start frontend with backend URL pointing to Minikube ingress
cd frontend
cat > .env.local <<'EOF'
NEXT_PUBLIC_API_URL=http://go-task-orchestrator.local/api
NEXT_TELEMETRY_DISABLED=1
EOF
npm ci
npm run dev
```

Access:

- Frontend UI: `http://localhost:3000`
- Backend API (via ingress): `http://go-task-orchestrator.local/api`

> Note: current Kubernetes manifests deploy backend services (gateway/auth/task). The frontend is not deployed in `k8s/` yet.

**Kubernetes Files:**
- `k8s/namespace.yaml` - Isolated namespace for the project
- `k8s/configmap.yaml` - Configuration (non-sensitive)
- `k8s/secret.yaml` - Secrets management (JWT_SECRET, credentials)
- `k8s/api-gateway.yaml` - API Gateway deployment & service
- `k8s/auth-service.yaml` - Auth Service deployment & service
- `k8s/task-service.yaml` - Task Service deployment & service
- `k8s/ingress.yaml` - Ingress for external access

## 🧪 Testing

### Run Go Tests

```bash
# Test all packages
go test -v -race -coverprofile=coverage.out ./...

# Test specific package
go test -v ./cmd/api-gateway/...

# Generate coverage report
go tool cover -html=coverage.out
```

### Test Frontend Build

```bash
cd frontend
npm ci
npm run lint
npm run build
```

### Integration Testing

Start the full stack and run manual tests:

```bash
# Start services
docker compose up

# In another terminal, run tests
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"test123"}'
```

## 🔄 CI/CD Pipeline

GitHub Actions automatically:
1. ✅ Runs Go tests (`go test ./...`)
2. ✅ Builds all Go binaries
3. ✅ Lints and builds Next.js frontend
4. ✅ Validates docker-compose.yml configuration
5. ✅ Runs on push to `main`/`develop`, PRs, and weekly schedule

**Workflow:** `.github/workflows/ci-cd.yml`

## 📝 Development Guide

### Adding a New Endpoint

1. **Define handler** in service (`cmd/service/main.go`)
2. **Register route** in API Gateway (`cmd/api-gateway/main.go`)
3. **Add types** if needed (`frontend/src/types/index.ts`)
4. **Create frontend component** (`frontend/src/components/`)
5. **Test** locally before pushing

### Code Style
- **Go:** Follow standard Go conventions (`gofmt`)
- **TypeScript/React:** Follow ESLint rules (`npm run lint`)
- **Commits:** Use conventional commits (feat:, fix:, docs:, etc.)

## 🚢 Deployment

### Deploy to Production

1. **Set up secrets** in GitHub (for Docker Hub credentials if pushing images)
2. **Push to main branch** - CI/CD runs automatically
3. **Verify tests pass** in GitHub Actions
4. **Deploy** using Kubernetes manifests or Docker Compose on your server

### Health Checks

All services expose health endpoints:
```bash
curl http://localhost:8080/health      # API Gateway
curl http://localhost:8081/health      # Auth Service
curl http://localhost:8082/health      # Task Service
```

## 🔐 Security Considerations

- ✅ JWT token-based authentication
- ✅ Password hashing (auth-service)
- ✅ Role-based access control (RBAC) for user listing
- ✅ Request validation on all endpoints
- ✅ CORS support for frontend communication
- ⚠️ **TODO:** HTTPS in production (configure in deployment)
- ⚠️ **TODO:** Database integration (currently in-memory)
- ⚠️ **TODO:** Rate limiting

## 📖 Additional Documentation

- [API Reference](docs/API_REFERENCE.md)
- [Database Schema](docs/DATABASE_SCHEMA.md)
- [Project Overview](docs/PROJECT_OVERVIEW.md)
- [Setup Guide](docs/SETUP_GUIDE.md)
- [Deployment Guide](docs/DEPLOYMENT.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'feat: add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

## 📞 Support

For issues, questions, or suggestions:
- Create a GitHub Issue
- Check existing documentation in `/docs`
- Review the API documentation in `API_DOCUMENTATION.md`

---

**Last Updated:** March 2024
**Status:** ✅ Production Ready
**Version:** 1.0.0
