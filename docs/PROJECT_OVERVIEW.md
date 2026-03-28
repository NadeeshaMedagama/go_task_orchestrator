# Project Overview

## Introduction

`Go_Task_Orchestrator` is a Go microservices project for authentication, authorization, and task orchestration.

The platform is organized around independent services behind one API gateway. This keeps service boundaries clear and enables independent deployment and scaling.

## Architecture

```text
Client
  |
  v
api-gateway (Go, :8080)
  |----------------------------|
  v                            v
auth-service (:8081)       task-service (:8082)
```

### Service Responsibilities

- `api-gateway`: public entrypoint, JWT validation, request routing
- `auth-service`: user register/login, token issuance, admin-only user listing
- `task-service`: task CRUD, status updates, filtering, pagination, ownership checks

## Technology Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22 |
| Auth | JWT (`github.com/golang-jwt/jwt/v5`) |
| HTTP | `net/http` (standard library) |
| Local orchestration | Docker Compose |
| Kubernetes | Native manifests + kustomize |
| CI | GitHub Actions (`.github/workflows/ci-cd.yml`) |

## Current Data Model

Current implementation uses in-memory stores inside `auth-service` and `task-service`.

- Good for local development and API contract validation
- Data resets when service restarts
- Persistent database integration can be added in a later phase

## Repository Structure

```text
Go_Task_Orchestrator/
├── cmd/
│   ├── api-gateway/
│   ├── auth-service/
│   └── task-service/
├── internal/platform/jwtutil/
├── deployments/docker/
├── docs/
├── k8s/
├── docker-compose.yml
├── go.mod
└── action.yml
```

## Security Model

- Public endpoints: `/api/auth/register`, `/api/auth/login`
- Protected endpoints: `/api/tasks*`, `/api/users`
- Gateway validates Bearer token and forwards user context via headers:
  - `X-User-Name`
  - `X-User-Role`
- `ADMIN` can list users and access all tasks
- `USER` can access only own tasks

## Related Documentation

- [Setup Guide](./SETUP_GUIDE.md)
- [API Reference](./API_REFERENCE.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [API Documentation](../API_DOCUMENTATION.md)
