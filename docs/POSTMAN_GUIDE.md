# Postman API Testing Guide - Go Task Orchestrator

## Overview
This project uses a microservices architecture with an API Gateway pattern. All requests should go through the API Gateway.

### ⚠️ Important: Deployment Method

- **Docker Compose (Local):** Use `http://localhost:8080`
- **Minikube Cluster:** Use `http://go-task-orchestrator.local` (requires `/etc/hosts` setup)
- **Minikube Port Forwarding:** Use `http://localhost:8080` (requires port-forward command)

---

## 🔧 Service Information

### Services & Ports
| Service | Port | URL | Purpose |
|---------|------|-----|---------|
| **API Gateway** | 8080 | `http://localhost:8080` | Main entry point, routes requests |
| **Auth Service** | 8081 | `http://localhost:8081` | User authentication & management |
| **Task Service** | 8082 | `http://localhost:8082` | Task management operations |
| **Frontend** | 3000 | `http://localhost:3000` | Next.js UI |

### Base URL
```
http://localhost:8080/api
```

---

## 🏥 Health Check Endpoints

Test service connectivity:

### API Gateway Health
```
GET http://localhost:8080/health
```

### Auth Service Health
```
GET http://localhost:8081/health
```

### Task Service Health
```
GET http://localhost:8082/health
```

---

## 🔐 Authentication Endpoints (Public)

### 1. Register User
```
POST http://localhost:8080/api/auth/register
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "john",
  "email": "john@example.com",
  "password": "secret123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "tokenType": "Bearer",
  "username": "john",
  "email": "john@example.com",
  "role": "USER"
}
```

---

### 2. Login User
```
POST http://localhost:8080/api/auth/login
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "john",
  "password": "secret123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "tokenType": "Bearer",
  "username": "john",
  "email": "john@example.com",
  "role": "USER"
}
```

**⚠️ Important:** Save the `token` value from the response. You'll need it for authenticated requests.

---

## 👥 User Endpoints

### List All Users (Admin Only)
```
GET http://localhost:8080/api/users
Authorization: Bearer <token>
```

**Required Header:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:**
```json
[
  {
    "username": "john",
    "email": "john@example.com",
    "role": "USER"
  },
  {
    "username": "admin",
    "email": "admin@example.com",
    "role": "ADMIN"
  }
]
```

---

## 📋 Task Endpoints (All Require Authentication)

### Required Header for All Task Endpoints:
```
Authorization: Bearer <token>
```

---

### 1. List Tasks (GET)
```
GET http://localhost:8080/api/tasks
Authorization: Bearer <token>
```

**Query Parameters:**
| Parameter | Type | Values | Default | Example |
|-----------|------|--------|---------|---------|
| `status` | string | `TODO`, `IN_PROGRESS`, `DONE` | - | `?status=TODO` |
| `priority` | string | `LOW`, `MEDIUM`, `HIGH` | - | `?priority=HIGH` |
| `page` | integer | Any number | `0` | `?page=0` |
| `size` | integer | Any number | `10` | `?size=10` |

**Example Requests:**
```
GET http://localhost:8080/api/tasks
GET http://localhost:8080/api/tasks?status=TODO
GET http://localhost:8080/api/tasks?priority=HIGH&page=0&size=10
GET http://localhost:8080/api/tasks?status=IN_PROGRESS&priority=MEDIUM
```

**Response:**
```json
{
  "content": [
    {
      "id": 1,
      "title": "Fix login bug",
      "description": "Investigate auth flow issue",
      "status": "TODO",
      "priority": "HIGH",
      "dueDate": "2026-05-01",
      "owner": "john",
      "createdAt": "2026-03-23T10:30:00Z",
      "updatedAt": "2026-03-23T10:30:00Z"
    }
  ],
  "page": 0,
  "size": 10,
  "totalElements": 1
}
```

---

### 2. Create Task (POST)
```
POST http://localhost:8080/api/tasks
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "title": "Fix login",
  "description": "Investigate auth flow",
  "status": "TODO",
  "priority": "HIGH",
  "dueDate": "2026-05-01"
}
```

**Status Values:** `TODO`, `IN_PROGRESS`, `DONE`
**Priority Values:** `LOW`, `MEDIUM`, `HIGH`

**Response:**
```json
{
  "id": 1,
  "title": "Fix login",
  "description": "Investigate auth flow",
  "status": "TODO",
  "priority": "HIGH",
  "dueDate": "2026-05-01",
  "owner": "john",
  "createdAt": "2026-03-23T10:30:00Z",
  "updatedAt": "2026-03-23T10:30:00Z"
}
```

---

### 3. Get Task by ID
```
GET http://localhost:8080/api/tasks/{id}
Authorization: Bearer <token>
```

**Example:**
```
GET http://localhost:8080/api/tasks/1
```

**Response:**
```json
{
  "id": 1,
  "title": "Fix login",
  "description": "Investigate auth flow",
  "status": "TODO",
  "priority": "HIGH",
  "dueDate": "2026-05-01",
  "owner": "john",
  "createdAt": "2026-03-23T10:30:00Z",
  "updatedAt": "2026-03-23T10:30:00Z"
}
```

---

### 4. Update Task (PUT)
```
PUT http://localhost:8080/api/tasks/{id}
Content-Type: application/json
Authorization: Bearer <token>
```

**Example:**
```
PUT http://localhost:8080/api/tasks/1
```

**Request Body:**
```json
{
  "title": "Fix login - Updated",
  "description": "Investigate auth flow - Complete refactor needed",
  "status": "IN_PROGRESS",
  "priority": "MEDIUM",
  "dueDate": "2026-06-01"
}
```

**Response:** Updated task object (same as GET)

---

### 5. Update Task Status (PATCH)
```
PATCH http://localhost:8080/api/tasks/{id}/status
Content-Type: application/json
Authorization: Bearer <token>
```

**Example:**
```
PATCH http://localhost:8080/api/tasks/1/status
```

**Request Body:**
```json
{
  "status": "DONE"
}
```

**Response:**
```json
{
  "id": 1,
  "title": "Fix login",
  "description": "Investigate auth flow",
  "status": "DONE",
  "priority": "HIGH",
  "dueDate": "2026-05-01",
  "owner": "john",
  "createdAt": "2026-03-23T10:30:00Z",
  "updatedAt": "2026-03-23T10:30:00Z"
}
```

---

### 6. Delete Task
```
DELETE http://localhost:8080/api/tasks/{id}
Authorization: Bearer <token>
```

**Example:**
```
DELETE http://localhost:8080/api/tasks/1
```

**Response:** HTTP 204 No Content (or success message)

---

## 🔑 JWT Token Management in Postman

### Method 1: Manual Token Usage
1. Login with credentials using `POST /auth/login`
2. Copy the `token` value from response
3. In your Postman request, go to **Authorization** tab
4. Select **Bearer Token** type
5. Paste the token in the **Token** field

### Method 2: Automatic Token Management (Recommended)
1. Create a **Postman Environment**
2. Add variable: `token` (initially empty)
3. In the login request, add a **Tests** script:
```javascript
if (pm.response.code === 200) {
    pm.environment.set("token", pm.response.json().token);
}
```
4. Use `{{token}}` in Authorization headers for other requests

---

## ⚙️ Environment Setup for Postman

### Create Environment Variables

Create a Postman Environment with these variables:

```
base_url = http://localhost:8080/api
token = <leave empty initially>
username = john
password = secret123
email = john@example.com
```

### Use Variables in Requests
```
{{base_url}}/tasks
Authorization: Bearer {{token}}
```

---

## 📝 Example Test Workflow

### 1. Register New User
```
POST {{base_url}}/auth/register
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "test123"
}
```

### 2. Login
```
POST {{base_url}}/auth/login
{
  "username": "testuser",
  "password": "test123"
}
```
Save the token from response.

### 3. Create Task
```
POST {{base_url}}/tasks
Authorization: Bearer {{token}}
{
  "title": "Test Task",
  "description": "Testing task creation",
  "status": "TODO",
  "priority": "HIGH",
  "dueDate": "2026-04-01"
}
```

### 4. Get All Tasks
```
GET {{base_url}}/tasks?status=TODO
Authorization: Bearer {{token}}
```

### 5. Update Task Status
```
PATCH {{base_url}}/tasks/1/status
Authorization: Bearer {{token}}
{
  "status": "IN_PROGRESS"
}
```

### 6. Delete Task
```
DELETE {{base_url}}/tasks/1
Authorization: Bearer {{token}}
```

---

## 🚨 Common Issues & Solutions

### Issue: "Unauthorized" or 401 Error
- **Solution:** Ensure your JWT token is valid and included in the Authorization header
- Check that token hasn't expired
- Re-login to get a fresh token

### Issue: 404 Not Found
- **Solution:** Verify the task ID exists
- Check the base URL is `http://localhost:8080/api`

### Issue: CORS Error
- **Solution:** This is expected in browser requests; Postman doesn't have CORS restrictions
- Frontend requests may need CORS headers configured

### Issue: Task Service Returns 500
- **Solution:** Ensure all services are running
- Check Docker containers are up: `docker ps`
- Check environment variables are set correctly

---

## 🐳 Docker Compose Setup

To run all services:

```bash
docker-compose up -d
```

To check services status:
```bash
docker-compose ps
```

To view logs:
```bash
docker-compose logs -f api-gateway
```

---

## 📚 Additional Resources

- **API Documentation:** See `API_DOCUMENTATION.md`
- **Database Schema:** See `docs/DATABASE_SCHEMA.md`
- **Deployment Guide:** See `docs/DEPLOYMENT.md`

---

**Last Updated:** March 23, 2026

