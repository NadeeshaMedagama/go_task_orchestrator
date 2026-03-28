# API Reference

Base URL: `http://localhost:8080/api`

## Authentication

### `POST /auth/register` (public)
Creates a user with default role `USER`.

### `POST /auth/login` (public)
Returns:

- `token`
- `tokenType` (`Bearer`)
- `username`
- `email`
- `role`

## Users

### `GET /users` (admin only)
Requires `Authorization: Bearer <token>` and caller role `ADMIN`.

## Tasks (authenticated)

All endpoints below require `Authorization: Bearer <token>`.

### `GET /tasks`
Query parameters:

- `status`: `TODO`, `IN_PROGRESS`, `DONE`
- `priority`: `LOW`, `MEDIUM`, `HIGH`
- `page`: default `0`
- `size`: default `10`

Returns paged response:

- `content`
- `page`
- `size`
- `totalElements`

### `POST /tasks`
Creates a task.

### `GET /tasks/{id}`
Fetches a task by ID.

### `PUT /tasks/{id}`
Replaces a task.

### `PATCH /tasks/{id}/status`
Updates only task status.

### `DELETE /tasks/{id}`
Deletes a task.

## Health Endpoints

- Gateway: `GET /health` on port `8080`
- Auth service: `GET /health` on port `8081`
- Task service: `GET /health` on port `8082`

## Error Conventions

Errors return standard HTTP status codes with JSON payloads like:

```json
{
  "error": "message"
}
```
