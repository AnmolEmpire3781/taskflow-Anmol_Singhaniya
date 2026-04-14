## 1. Overview

Backend-only implementation of **TaskFlow**, a task management system with JWT authentication, PostgreSQL, relational project/task ownership, automatic schema migrations, seed data, Docker Compose startup, and a Postman collection. The app is intended to run with `docker compose up`, uses migrations, hashes passwords with bcrypt, reads JWT secret from environment variables.

### Tech stack

- Go 1.24
- Chi router
- PostgreSQL 16
- pgx/pgxpool
- JWT (`github.com/golang-jwt/jwt/v5`)
- bcrypt (`golang.org/x/crypto/bcrypt`)
- Docker Compose
- Postman collection for API verification

### Core features

- Register and login
- JWT-protected APIs
- Projects CRUD
- Tasks CRUD
- Project task filtering by status and assignee
- Project stats endpoint
- Automatic migrations on startup
- Automatic seed loading on startup
- Structured JSON responses
- Structured logging
- Graceful shutdown

### Bonus features implemented

- Pagination on `GET /projects` and `GET /projects/:id/tasks`
- `GET /projects/:id/stats`

## 2. Architecture Decisions

### Why this structure

The backend is split into small packages to keep responsibilities separate:

- `config`: environment-driven configuration
- `auth`: password hashing and JWT generation/parsing
- `db`: connection, migration runner, seed runner
- `repository`: database access logic
- `service`: business rules and authorization checks
- `app`: router and application bootstrap
- `httpx`: reusable JSON response helpers
- `middleware`: auth middleware

This keeps handlers thin and makes code review easier.

### Tradeoffs

- I used straightforward SQL with `pgx` instead of an ORM to align with the requirement to manage schema via migrations instead of auto-migrate magic.
- Migrations are executed automatically at startup by the application using migration files from `backend/migrations`. That keeps `docker compose up` to a single command.
- I included a Postman collection instead of a frontend because backend-only candidates do not need the frontend, per the assignment spec.

### What I intentionally left out

- WebSocket/SSE live updates
- advanced RBAC
- refresh tokens
- full observability stack
- exhaustive integration test suite with ephemeral test DB

Those are improvements I would add next after the core submission is fully stable.

## 3. Running Locally

### Prerequisites

--Assume only **Docker Desktop / Docker Engine with Compose** is installed.

### Commands

```bash
git clone https://github.com/AnmolEmpire3781/taskflow_Anmol_Singhaniya
cd taskflow-Anmol_Singhaiya
cp .env.example .env
docker compose up --build
```

### App URL

- API: `http://localhost:8080`
- Health: `http://localhost:8080/health`

The backend container automatically:

1. waits for Postgres to be healthy,
2. connects to the database,
3. runs all `.up.sql` migrations,
4. applies seed data,
5. starts the HTTP server.

## 4. Running Migrations

Migrations run automatically on application startup.

Migration files are stored in:

```text
backend/migrations/
```

Both up and down migration files are included for each migration, as required by the assignment.

If you want to run migrations manually in the future, you can plug in `golang-migrate`, but for the assignment a no-manual-step startup flow is already provided.

## 5. Test Credentials

Use this seeded account right after startup:

```text
Email:    test@example.com
Password: password123
```

Seed data also includes:

- 1 project
- 3 tasks in different statuses (`todo`, `in_progress`, `done`)

This matches the seed requirement from the assignment.

## 6. API Reference

### Authentication

#### POST `/auth/register`

Request:

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "secret123"
}
```

#### POST `/auth/login`

Request:

```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "token": "<jwt>",
  "user": {
    "id": "uuid",
    "name": "Test User",
    "email": "test@example.com",
    "created_at": "2026-04-13T00:00:00Z"
  }
}
```

### Projects

- `GET /projects?page=&limit=`
- `POST /projects`
- `GET /projects/:id`
- `PATCH /projects/:id`
- `DELETE /projects/:id`
- `GET /projects/:id/stats`

### Tasks

- `GET /projects/:id/tasks?status=&assignee=&page=&limit=`
- `POST /projects/:id/tasks`
- `PATCH /tasks/:id`
- `DELETE /tasks/:id`

### Error format

Validation error:

```json
{
  "error": "validation failed",
  "fields": {
    "email": "is required"
  }
}
```

Unauthenticated:

```json
{ "error": "unauthorized" }
```

Forbidden:

```json
{ "error": "forbidden" }
```

Not found:

```json
{ "error": "not found" }
```

### Postman collection

See:

```text
postman/taskflow.postman_collection.json
```

## 7. What I’d Do With More Time

- Add table-driven integration tests using a disposable Postgres test container
- Add request logging middleware with latency and status fields
- Add refresh tokens and token revocation
- Would add API integration tests for auth, project, and task flows.

<img width="1352" height="949" alt="image" src="https://github.com/user-attachments/assets/5e2bec89-ea5f-4016-a63f-f0a577d284aa" />


<img width="1919" height="703" alt="image" src="https://github.com/user-attachments/assets/f0afc19a-15cb-4a67-9541-1a47fde7b2c9" />

<img width="1919" height="992" alt="image" src="https://github.com/user-attachments/assets/57141009-28e1-48d8-9644-e820c3f905d9" />

<img width="1919" height="990" alt="image" src="https://github.com/user-attachments/assets/63f242c5-11f4-418c-9e02-4a3e5f0a17ef" />

<img width="1918" height="1006" alt="image" src="https://github.com/user-attachments/assets/27516127-2b24-49be-9e7e-af39d0a2f665" />

