# Secure Chat Service

A real-time, secure chat application built with Go, Gin, PostgreSQL, Redis, NATS, and WebSockets. The project is organized around a hexagonal architecture with clear separation between domain logic, application services, and infrastructure adapters.

## Overview

This service provides:

- User registration and authentication
- Chat room creation and management
- Real-time messaging through WebSockets
- Message persistence with PostgreSQL
- Cached data access with Redis
- Event-driven messaging with NATS
- Basic security protections such as CSRF protection, rate limiting, and security headers

## Tech Stack

- Go
- Gin Web Framework
- GORM + PostgreSQL
- Redis
- NATS JetStream
- Gorilla WebSocket
- JWT authentication
- Docker / Docker Compose

## Project Structure

- cmd/server: application entrypoint
- internal/configs: environment and configuration loading
- internal/core: domain models and application services
- internal/adapters: HTTP handlers, WebSocket handlers, repositories, and external integrations
- pkg: shared helpers such as logging, hashing, and encryption
- templates: HTML templates for the web UI

## Features

- Secure authentication with JWT
- Session-based CSRF protection
- Rate limiting and security headers
- WebSocket-based chat communication
- Persistent message storage
- Encrypted Redis cache support

## Prerequisites

Before running the project, make sure you have:

- Go 1.26 or newer
- Docker and Docker Compose
- A working terminal environment

## Environment Variables

Create a .env file in the project root with values similar to:

```env
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=secure_chat_db
REDIS_URL=redis://localhost:6379/0
NATS_URL=nats://localhost:4222
JWT_SECRET=your-jwt-secret
CSRF_SECRET=your-csrf-secret
CACHE_ENCRYPTION_KEY=64-character-hex-string
```

> The application expects these values to be available when it starts.

## Running with Docker Compose

The easiest way to start the full stack is with Docker Compose:

```bash
docker compose up --build
```

This will start:

- PostgreSQL
- Redis
- NATS
- The chat service on port 8080

Open the app in your browser at:

```text
http://localhost:8080
```

## Running Locally

If you want to run the server directly on your machine:

1. Start the supporting services:

```bash
docker compose up redis postgres nats
```

2. Run the application:

```bash
go run ./cmd/server
```

# Contributation
