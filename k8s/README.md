# Kubernetes Setup (Minikube)

This directory contains Kubernetes manifests for `Go_Task_Orchestrator`.

## Included Manifests

- `namespace.yaml`
- `configmap.yaml`
- `secret.yaml`
- `auth-service.yaml`
- `task-service.yaml`
- `api-gateway.yaml`
- `frontend.yaml`
- `ingress.yaml`
- `kustomization.yaml`

## 1) Start Minikube and ingress

```bash
minikube start
minikube addons enable ingress
```

## 2) Build local images inside Minikube

```bash
eval "$(minikube docker-env)"
docker build -t go-task-orchestrator/auth-service:local -f deployments/docker/auth-service.Dockerfile .
docker build -t go-task-orchestrator/task-service:local -f deployments/docker/task-service.Dockerfile .
docker build -t go-task-orchestrator/api-gateway:local -f deployments/docker/api-gateway.Dockerfile .
docker build -t go-task-orchestrator/frontend:local -f frontend/Dockerfile frontend
```

## 3) Update secret

Edit `k8s/secret.yaml` and set a strong `JWT_SECRET` value.

## 4) Deploy

```bash
kubectl apply -k k8s
```

## 5) Verify pods

```bash
kubectl get pods -n go-task-orchestrator -w
```

## 6) Configure local host

```bash
minikube ip
```

Add to `/etc/hosts`:

```text
<MINIKUBE_IP> go-task-orchestrator.local
```

Then test:

- Frontend: `http://go-task-orchestrator.local/`
- Gateway health: `http://go-task-orchestrator.local/health`
- API base: `http://go-task-orchestrator.local/api`

## Troubleshooting

```bash
kubectl get all -n go-task-orchestrator
kubectl logs -n go-task-orchestrator deploy/api-gateway
kubectl logs -n go-task-orchestrator deploy/auth-service
kubectl logs -n go-task-orchestrator deploy/task-service
```
