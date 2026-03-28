# ✅ Your Backend is Working! - Connection Guide

Your Go Task Orchestrator backend is **successfully running** in Minikube. The error you got (`ECONNREFUSED 127.0.0.1:8080`) happened because the services in Minikube are not accessible at `localhost:8080` directly. Here's how to fix it:

---

## 🚀 Quick Start (5 minutes)

### Option A: Using Minikube Ingress (Recommended)

**Step 1:** Add hostname to your `/etc/hosts` file
```bash
echo "192.168.49.2 go-task-orchestrator.local" | sudo tee -a /etc/hosts
```

**Step 2:** Verify it worked
```bash
ping go-task-orchestrator.local
```

**Step 3:** Update Postman Base URL
```
http://go-task-orchestrator.local/api
```

✅ **You're done!** Test with this health check:
```
GET http://go-task-orchestrator.local/health
```

---

### Option B: Using Port Forwarding

**Step 1:** In your terminal, run:
```bash
kubectl port-forward -n go-task-orchestrator svc/api-gateway 8080:8080
```

You'll see:
```
Forwarding from 127.0.0.1:8080 -> 8080
Listening on port 8080
```

**Step 2:** Update Postman Base URL
```
http://localhost:8080/api
```

✅ **Keep the terminal running!** The port forward must stay active.

---

## ✅ Verify Your Services are Running

### Check Pod Status
```bash
kubectl get pods -n go-task-orchestrator
```

Expected:
```
NAME                           READY   STATUS    
api-gateway-6b6d48c548-4wchk   1/1     Running   
auth-service-5c6f88d4f-djbsk   1/1     Running   
task-service-dbd49c6d7-bdp2v   1/1     Running   
```

### Check Services
```bash
kubectl get svc -n go-task-orchestrator
```

Expected:
```
NAME           TYPE        CLUSTER-IP       PORT(S)
api-gateway    ClusterIP   10.106.196.181   8080/TCP
auth-service   ClusterIP   10.104.52.7      8081/TCP
task-service   ClusterIP   10.101.35.76     8082/TCP
```

### Check Ingress
```bash
kubectl get ingress -n go-task-orchestrator
```

Expected:
```
NAME                           CLASS   HOSTS                        ADDRESS      PORTS
go-task-orchestrator-ingress   nginx   go-task-orchestrator.local   192.168.49.2  80
```

---

## 🧪 Test Connectivity

### Using curl with Port Forwarding

Terminal 1:
```bash
kubectl port-forward -n go-task-orchestrator svc/api-gateway 8080:8080
```

Terminal 2:
```bash
curl http://localhost:8080/health
```

Response:
```json
{"status":"ok","service":"api-gateway"}
```

### Using curl with Ingress

```bash
curl http://go-task-orchestrator.local/health
```

---

## 📊 Service Overview

| Component | Status | Access Method 1 | Access Method 2 |
|-----------|--------|---|---|
| API Gateway | ✅ Running | `go-task-orchestrator.local:80` | `localhost:8080` (port-forward) |
| Auth Service | ✅ Running | Internal (via gateway) | `localhost:8081` (port-forward) |
| Task Service | ✅ Running | Internal (via gateway) | `localhost:8082` (port-forward) |

---

## 📚 Next Steps for Postman Testing

Once you have either **Option A** or **Option B** set up:

1. **Base URL in Postman:** Use the appropriate URL above
2. **Register a user:**
   ```
   POST /api/auth/register
   {
     "username": "testuser",
     "email": "test@example.com", 
     "password": "test123"
   }
   ```
3. **Login:**
   ```
   POST /api/auth/login
   {
     "username": "testuser",
     "password": "test123"
   }
   ```
4. **Save the token** from login response
5. **Create tasks** using the token in Authorization header

See `POSTMAN_GUIDE.md` and `MINIKUBE_TESTING_GUIDE.md` for complete API documentation.

---

## 🛠️ Helpful Commands

```bash
# View all resources
kubectl get all -n go-task-orchestrator

# View pod logs
kubectl logs -n go-task-orchestrator api-gateway-6b6d48c548-4wchk -f

# Port forward API gateway
kubectl port-forward -n go-task-orchestrator svc/api-gateway 8080:8080

# Port forward all services
kubectl port-forward -n go-task-orchestrator svc/auth-service 8081:8081 &
kubectl port-forward -n go-task-orchestrator svc/task-service 8082:8082 &

# Check Minikube IP
minikube ip

# Run health check script
./scripts/health-check.sh
```

---

## 🎯 Why This Happened

1. **Services in Minikube** run inside a virtual cluster, not on your host machine
2. **Direct localhost access** doesn't work without port forwarding
3. **Ingress** provides a DNS-based way to access services (like a real domain)

---

## 📞 Troubleshooting

### Still can't connect?

1. **Check pods are really running:**
   ```bash
   kubectl get pods -n go-task-orchestrator
   ```

2. **Check pod logs for errors:**
   ```bash
   kubectl logs -n go-task-orchestrator api-gateway-6b6d48c548-4wchk
   ```

3. **Test directly in pod:**
   ```bash
   kubectl exec -it -n go-task-orchestrator api-gateway-6b6d48c548-4wchk -- curl http://localhost:8080/health
   ```

4. **Check if /etc/hosts was updated correctly:**
   ```bash
   cat /etc/hosts | grep go-task-orchestrator
   ```

5. **Flush DNS cache (macOS):**
   ```bash
   sudo dscacheutil -flushcache
   ```

---

## 📖 Documentation Files

- **POSTMAN_GUIDE.md** - Complete API endpoint reference
- **MINIKUBE_TESTING_GUIDE.md** - Detailed Minikube setup guide
- **scripts/health-check.sh** - Automated health check script

---

**Status:** ✅ Backend is healthy and ready to test!

Choose your preferred connection method above and start testing with Postman! 🚀

