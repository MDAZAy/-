# Auto Service Booking System

Web application for online booking in an auto service.

## Stack

- Frontend: Blazor WebAssembly + MudBlazor
- Backend: Go + Gin + GORM
- Database: MySQL 8
- Auth: JWT + Refresh Token
- Docs: Swagger
- Migrations: golang-migrate

## Structure

- `backend/`
- `frontend/`
- `docker-compose.yml`

## Run

1. Copy `backend/.env.example` to `backend/.env`
2. Start MySQL with `docker compose up -d mysql`
3. Run backend with `go run ./cmd/api` from `backend/`
4. Run frontend with `dotnet run` from `frontend/`

## Seed users

- `admin@example.com / Admin123`
- `user@example.com / User123`
