# Quick Reference Card - Minikube Backend Testing

## 🎯 Your Current Setup

```
Minikube IP: 192.168.49.2
Services: api-gateway, auth-service, task-service (ALL RUNNING ✅)
Ingress: go-task-orchestrator.local
```

---

## 🚀 Get Started (Pick One)

### Option 1️⃣: Ingress (Recommended - Works like production)

```bash
# 1. One-time setup (add to /etc/hosts)
echo "192.168.49.2 go-task-orchestrator.local" | sudo tee -a /etc/hosts

# 2. Postman Base URL
http://go-task-orchestrator.local/api

# 3. Test
curl http://go-task-orchestrator.local/health
```

---

### Option 2️⃣: Port Forwarding (Simple setup)

```bash
# 1. In Terminal, run:
kubectl port-forward -n go-task-orchestrator svc/api-gateway 8080:8080

# 2. Postman Base URL
http://localhost:8080/api

# 3. Test (in another terminal)
curl http://localhost:8080/health
```

---

## ✅ Verify Services

```bash
# All pods running?
kubectl get pods -n go-task-orchestrator

# Services created?
kubectl get svc -n go-task-orchestrator

# Ingress configured?
kubectl get ingress -n go-task-orchestrator

# Check logs
kubectl logs -n go-task-orchestrator api-gateway-6b6d48c548-4wchk -f
```

---

## 🔑 Postman Requests

### Health Check
```
GET http://go-task-orchestrator.local/health
```

### Register
```
POST http://go-task-orchestrator.local/api/auth/register
Content-Type: application/json

{
  "username": "john",
  "email": "john@example.com",
  "password": "secret123"
}
```

### Login (Save the token!)
```
POST http://go-task-orchestrator.local/api/auth/login
Content-Type: application/json

{
  "username": "john",
  "password": "secret123"
}
```

### Create Task
```
POST http://go-task-orchestrator.local/api/tasks
Content-Type: application/json
Authorization: Bearer <TOKEN_FROM_LOGIN>

{
  "title": "My Task",
  "description": "Task description",
  "status": "TODO",
  "priority": "HIGH",
  "dueDate": "2026-04-01"
}
```

### Get Tasks
```
GET http://go-task-orchestrator.local/api/tasks
Authorization: Bearer <TOKEN>
```

### Update Status
```
PATCH http://go-task-orchestrator.local/api/tasks/1/status
Authorization: Bearer <TOKEN>

{
  "status": "DONE"
}
```

---

## 🐛 If It's Not Working

| Problem | Solution |
|---------|----------|
| `ECONNREFUSED 127.0.0.1:8080` | Use port-forward or Ingress URL, not localhost |
| DNS not resolving `go-task-orchestrator.local` | Run: `sudo dscacheutil -flushcache` (macOS) |
| Pod not running | Check: `kubectl logs -n go-task-orchestrator <pod-name>` |
| Port already in use | Kill existing process: `lsof -i :8080` then `kill -9 <PID>` |

---

## 📚 More Info

- **POSTMAN_GUIDE.md** - Full API reference
- **MINIKUBE_TESTING_GUIDE.md** - Detailed setup
- **scripts/health-check.sh** - Run: `./scripts/health-check.sh`

---

**Status: ✅ BACKEND IS READY! Pick Option 1 or 2 above and start testing!**

