# Data Model Reference

## Current State

`Go_Task_Orchestrator` currently stores users and tasks in memory inside service processes:

- `auth-service`: user accounts and roles
- `task-service`: task entities

This means data is not persisted across restarts.

## Logical Entities

### User

| Field | Type | Notes |
|---|---|---|
| `username` | string | unique key in auth-service |
| `email` | string | user email |
| `password` | string | plain in current local-only demo implementation |
| `role` | string | `USER` or `ADMIN` |

### Task

| Field | Type | Notes |
|---|---|---|
| `id` | int64 | auto-increment in task-service |
| `title` | string | required |
| `description` | string | optional |
| `status` | string | `TODO`, `IN_PROGRESS`, `DONE` |
| `priority` | string | `LOW`, `MEDIUM`, `HIGH` |
| `dueDate` | string | optional |
| `owner` | string | username of owner |
| `createdAt` | timestamp | set on create |
| `updatedAt` | timestamp | set on update |

## Planned Persistence (Next Phase)

For production readiness, move from in-memory state to a database-backed architecture.

Recommended minimum:

- PostgreSQL
- Separate schemas/tables per service boundary where practical
- Password hashing in `auth-service`
- Migration tooling (for example, `golang-migrate`)

## Notes

This document intentionally reflects the current code behavior in this repository.
