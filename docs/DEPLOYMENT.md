# Deployment & CI/CD

This document describes deployment for `Go_Task_Orchestrator`.

## Docker

Three Dockerfiles are used:

- `deployments/docker/api-gateway.Dockerfile`
- `deployments/docker/auth-service.Dockerfile`
- `deployments/docker/task-service.Dockerfile`

Each image is multi-stage:

1. Build Go binary in `golang:1.22-alpine`
2. Run binary in `alpine:3.20`

## Docker Compose

`docker-compose.yml` starts:

- `api-gateway` on `:8080`
- `auth-service` on `:8081`
- `task-service` on `:8082`

Run:

```bash
docker compose up --build
```

## CI Workflow

`/.github/workflows/ci-cd.yml` runs on push and pull request.

Pipeline steps:

1. `go mod download`
2. `go test ./...`
3. build all service binaries
4. `docker compose config`

## Marketplace Action

`action.yml` provides a composite action that deploys the stack using Docker Compose.

Required input:

- `jwt-secret`

Optional input:

- `gateway-port` (default `8080`)

## Kubernetes Deployment

Kubernetes manifests are in `k8s/` and deploy the same microservices topology.

Apply:

```bash
kubectl apply -k k8s
```

## Production Notes

- Replace `JWT_SECRET` with a strong secret
- Use private image registry and immutable tags
- Configure ingress TLS for production domains
- Add persistent storage and a real database before production usage
