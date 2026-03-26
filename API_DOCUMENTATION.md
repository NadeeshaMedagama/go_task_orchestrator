# Go_Task_Orchestrator API Documentation

Base URL: `http://localhost:8080/api`

## Authentication

### `POST /auth/register`
Create a user.

Request:
```json
{
  "username": "john",
  "email": "john@example.com",
  "password": "secret123"
}
```

### `POST /auth/login`
Login and receive JWT token.

Request:
```json
{
  "username": "john",
  "password": "secret123"
}
```

Response:
```json
{
  "token": "<jwt>",
  "tokenType": "Bearer",
  "username": "john",
  "email": "john@example.com",
  "role": "USER"
}
```

### `GET /users`
List users (admin only).

Required header:
- `Authorization: Bearer <token>`

## Tasks

All task endpoints require:
- `Authorization: Bearer <token>`

### `GET /tasks`
Query params:
- `status`: `TODO`, `IN_PROGRESS`, `DONE`
- `priority`: `LOW`, `MEDIUM`, `HIGH`
- `page`: default `0`
- `size`: default `10`

### `POST /tasks`
Request:
```json
{
  "title": "Fix login",
  "description": "Investigate auth flow",
  "status": "TODO",
  "priority": "HIGH",
  "dueDate": "2026-05-01"
}
```

### `GET /tasks/{id}`
Fetch a task by ID.

### `PUT /tasks/{id}`
Replace a task.

### `PATCH /tasks/{id}/status`
Update only status.

Request:
```json
{
  "status": "DONE"
}
```

### `DELETE /tasks/{id}`
Delete task.
