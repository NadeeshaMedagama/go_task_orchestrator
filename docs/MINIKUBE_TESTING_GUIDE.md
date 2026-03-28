# Minikube Testing Guide - Go Task Orchestrator

## Current Minikube Setup

Your services are running in a Minikube cluster with the following configuration:

```
Minikube IP: 192.168.49.2
Namespace: go-task-orchestrator
Ingress Host: go-task-orchestrator.local
Ingress Port: 80 (nginx)
```

### Running Pods
```
api-gateway-6b6d48c548-4wchk       Running  Port 8080
auth-service-5c6f88d4f-djbsk       Running  Port 8081
task-service-dbd49c6d7-bdp2v       Running  Port 8082
```

---

## ✅ Method 1: Using Ingress (Recommended)

### Step 1: Update Your `/etc/hosts` File

Add the following line to your `/etc/hosts` file:

```bash
sudo nano /etc/hosts
```

Add this line:
```
192.168.49.2 go-task-orchestrator.local
```

Or use this command:
```bash
echo "192.168.49.2 go-task-orchestrator.local" | sudo tee -a /etc/hosts
```

### Step 2: Update Postman Base URL

Use this base URL in Postman:
```
http://go-task-orchestrator.local/api
```

### Step 3: Test Health Endpoint

```
GET http://go-task-orchestrator.local/health
```

Should return:
```json
{
  "status": "ok",
  "service": "api-gateway"
}
```

---

## ✅ Method 2: Port Forwarding (Alternative)

If you prefer not to modify `/etc/hosts`, use port forwarding:

### Forward API Gateway Port

```bash
kubectl port-forward -n go-task-orchestrator svc/api-gateway 8080:8080
```

This will output:
```
Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
Listening on port 8080
```

### Use Postman with Localhost

Base URL:
```
http://localhost:8080/api
```

Now you can use the standard Postman guide!

---

## ✅ Method 3: Direct Service Access

Access services directly in Minikube:

```bash
# API Gateway
kubectl port-forward -n go-task-orchestrator svc/api-gateway 8080:8080

# Auth Service
kubectl port-forward -n go-task-orchestrator svc/auth-service 8081:8081

# Task Service
kubectl port-forward -n go-task-orchestrator svc/task-service 8082:8082
```

---

## 🔍 How to Check if Services are Working

### 1. Check Pod Status
```bash
kubectl get pods -n go-task-orchestrator
```

Expected output: All pods in `Running` status with `1/1` Ready count.

### 2. Check Service Status
```bash
kubectl get svc -n go-task-orchestrator
```

Expected output: All services should show `ClusterIP` type with their respective ports.

### 3. Check Pod Logs

**API Gateway logs:**
```bash
kubectl logs -n go-task-orchestrator api-gateway-6b6d48c548-4wchk -f
```

**Auth Service logs:**
```bash
kubectl logs -n go-task-orchestrator auth-service-5c6f88d4f-djbsk -f
```

**Task Service logs:**
```bash
kubectl logs -n go-task-orchestrator task-service-dbd49c6d7-bdp2v -f
```

### 4. Test with curl from Inside Cluster

```bash
# Get into the api-gateway pod
kubectl exec -it -n go-task-orchestrator api-gateway-6b6d48c548-4wchk -- sh

# Inside the pod, test the health endpoint
curl http://localhost:8080/health
```

### 5. Test with Port Forwarding + curl

In one terminal, set up port forwarding:
```bash
kubectl port-forward -n go-task-orchestrator svc/api-gateway 8080:8080
```

In another terminal, test:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok","service":"api-gateway"}
```

---

## 🐳 Complete Testing Workflow with Minikube

### Setup (One-time)

```bash
# 1. Update /etc/hosts
echo "192.168.49.2 go-task-orchestrator.local" | sudo tee -a /etc/hosts

# 2. Verify services
kubectl get pods -n go-task-orchestrator
kubectl get svc -n go-task-orchestrator

# 3. Check ingress
kubectl get ingress -n go-task-orchestrator
```

### Testing with Postman

**Base URL:** `http://go-task-orchestrator.local/api`

**Test 1: Health Check**
```
GET http://go-task-orchestrator.local/health
```

**Test 2: Register User**
```
POST http://go-task-orchestrator.local/api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "test123"
}
```

**Test 3: Login**
```
POST http://go-task-orchestrator.local/api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "test123"
}
```

**Test 4: Create Task**
```
POST http://go-task-orchestrator.local/api/tasks
Content-Type: application/json
Authorization: Bearer <token_from_login>

{
  "title": "Test Task",
  "description": "Testing from Minikube",
  "status": "TODO",
  "priority": "HIGH",
  "dueDate": "2026-04-01"
}
```

---

## 🛠️ Debugging Tips

### Issue: DNS Not Resolving `go-task-orchestrator.local`

**Solution 1:** Verify `/etc/hosts` entry:
```bash
cat /etc/hosts | grep go-task-orchestrator
```

**Solution 2:** Flush DNS cache:
```bash
# macOS
sudo dscacheutil -flushcache

# Linux
sudo systemctl restart nscd
```

### Issue: Connection Refused on Port 8080

**Solution:** The service is not accessible directly. You need either:
1. Use Ingress hostname: `go-task-orchestrator.local`
2. Use port forwarding
3. Access from inside the cluster

### Issue: Pod Crashes or Errors

**Check logs:**
```bash
kubectl logs -n go-task-orchestrator <pod-name> --previous
```

**Describe pod:**
```bash
kubectl describe pod -n go-task-orchestrator <pod-name>
```

### Issue: Ingress Not Working

**Check Ingress Controller:**
```bash
kubectl get pods -n ingress-nginx
kubectl logs -n ingress-nginx <controller-pod>
```

**Check Ingress Configuration:**
```bash
kubectl describe ingress -n go-task-orchestrator go-task-orchestrator-ingress
```

---

## 📋 Quick Command Reference

```bash
# View all resources
kubectl get all -n go-task-orchestrator

# View pod details
kubectl get pods -n go-task-orchestrator -o wide

# Follow logs
kubectl logs -n go-task-orchestrator <pod-name> -f

# Execute command in pod
kubectl exec -it -n go-task-orchestrator <pod-name> -- <command>

# Port forward
kubectl port-forward -n go-task-orchestrator svc/<service-name> <local-port>:<service-port>

# Check service details
kubectl get svc -n go-task-orchestrator -o wide

# Check ingress
kubectl get ingress -n go-task-orchestrator -o wide

# Delete and redeploy
kubectl delete -f k8s/
kubectl apply -f k8s/
```

---

## 🚀 Recommended Testing Order

1. ✅ Check pod status: `kubectl get pods -n go-task-orchestrator`
2. ✅ Check services: `kubectl get svc -n go-task-orchestrator`
3. ✅ Test health endpoint via curl or Postman
4. ✅ Register a new user
5. ✅ Login and get token
6. ✅ Create a task
7. ✅ List tasks with filters
8. ✅ Update task status
9. ✅ Delete task

---

## 📚 Environment Configuration

For Postman, set up these variables in your Postman Environment:

```
base_url = http://go-task-orchestrator.local/api
token = <leave empty, will be set by login>
username = testuser
password = test123
email = test@example.com
```

Or if using port forwarding:

```
base_url = http://localhost:8080/api
token = <leave empty, will be set by login>
username = testuser
password = test123
email = test@example.com
```

---

**Last Updated:** March 23, 2026

