# Setup Guide

This guide explains how to run `Go_Task_Orchestrator` locally.

## Prerequisites

### Option A: Docker (recommended)

- Docker Engine 20.10+
- Docker Compose v2+

### Option B: Local Go runtime

- Go 1.22+

## Option A: Run with Docker Compose

1. Clone and enter project:

```bash
git clone https://github.com/<OWNER>/Go_Task_Orchestrator.git
cd Go_Task_Orchestrator
```

2. Create `.env` from template and set `JWT_SECRET`.

3. Start services:

```bash
docker compose up --build
```

4. Verify:

```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
```

## Option B: Run services directly

Open three terminals from the repository root.

Terminal 1:

```bash
export JWT_SECRET="$(openssl rand -base64 32)"
go run ./cmd/auth-service
```

Terminal 2:

```bash
go run ./cmd/task-service
```

Terminal 3:

```bash
export JWT_SECRET="<same-secret-used-for-auth-service>"
go run ./cmd/api-gateway
```

## Quick API Smoke Test

1. Register user:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret123"}'
```

2. Login:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret123"}'
```

3. Use token to create task:

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"First Task","status":"TODO","priority":"HIGH"}'
```

## Default Admin User

`auth-service` seeds a local admin user:

- Username: `admin`
- Password: `admin123`
- Role: `ADMIN`

## Environment Variables

| Variable | Default | Used by |
|---|---|---|
| `JWT_SECRET` | `change-me-in-production` | api-gateway, auth-service |
| `AUTH_SERVICE_URL` | `http://localhost:8081` | api-gateway |
| `TASK_SERVICE_URL` | `http://localhost:8082` | api-gateway |
| `PORT` | service-specific | all services |
