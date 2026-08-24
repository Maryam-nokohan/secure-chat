# Secure Chat Service 🔐

A real-time secure messaging platform built with **Go**, **WebSockets**, and an event-driven architecture.  
The project follows **Hexagonal Architecture (Ports & Adapters)** to keep business logic independent from infrastructure concerns.

Secure Chat provides scalable real-time communication with persistent storage, caching, asynchronous events, and security-focused features.

---

## ✨ Features

### Authentication & Security
- JWT-based authentication
- Password hashing
- CSRF protection
- Rate limiting
- Secure HTTP headers
- Encrypted cache storage

### Real-Time Communication
- WebSocket-based messaging
- Real-time chat rooms
- Event-driven message processing
- Asynchronous communication with NATS JetStream

### Data Management
- PostgreSQL for persistent storage
- Redis for caching and fast data access
- Message history storage
- Repository abstraction layer

### Architecture
- Hexagonal Architecture
- Clean separation of:
  - Domain logic
  - Application services
  - Infrastructure adapters
- Dependency inversion principles

---

# 🛠 Tech Stack

| Technology | Purpose |
|------------|---------|
| Go | Backend language |
| Gin | HTTP framework |
| PostgreSQL | Primary database |
| GORM | ORM |
| Redis | Cache layer |
| NATS JetStream | Event streaming |
| Gorilla WebSocket | Real-time communication |
| JWT | Authentication |
| Docker | Containerization |
| Docker Compose | Local development environment |

---

# 📁 Project Structure

```
├── cmd
│   └── server
│       └── main.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── primary
│   │   │   ├── http
│   │   │   │   ├── dto
│   │   │   │   │   ├── request.go
│   │   │   │   │   └── response.go
│   │   │   │   ├── handlers
│   │   │   │   │   ├── auth.go
│   │   │   │   │   ├── room.go
│   │   │   │   │   └── user.go
│   │   │   │   ├── middlewares
│   │   │   │   │   ├── auth.go
│   │   │   │   │   ├── csrf.go
│   │   │   │   │   ├── rate.go
│   │   │   │   │   └── security.go
│   │   │   │   └── routes
│   │   │   │       └── routes.go
│   │   │   └── websocket
│   │   │       ├── client.go
│   │   │       ├── handler.go
│   │   │       ├── hub.go
│   │   │       └── message.go
│   │   └── secondary
│   │       ├── auth
│   │       │   └── jwt_service.go
│   │       ├── nats
│   │       │   └── client.go
│   │       ├── postgres
│   │       │   ├── chat_repository.go
│   │       │   ├── db.go
│   │       │   ├── message_repository.go
│   │       │   ├── migrations
│   │       │   │   ├── migrate.go
│   │       │   │   └── schema.sql
│   │       │   └── user_repository.go
│   │       └── redis
│   │           ├── cache.go
│   │           ├── client.go
│   │           └── encryptedCache.go
│   ├── configs
│   │   ├── config.go
│   │   └── load.go
│   └── core
│       ├── application
│       │   ├── chat
│       │   │   ├── chat_service.go
│       │   │   └── pubsub.go
│       │   ├── message
│       │   │   └── message_service.go
│       │   └── user
│       │       └── user-service.go
│       ├── domain
│       │   ├── auth
│       │   │   └── auth.go
│       │   ├── chat
│       │   │   ├── entity.go
│       │   │   └── event.go
│       │   ├── message
│       │   │   ├── entity.go
│       │   │   └── pubsubmessage.go
│       │   └── user
│       │       └── entity.go
│       └── ports
│           ├── broker.go
│           ├── cache.go
│           ├── chat.go
│           ├── chat_repository.go
│           ├── message.go
│           ├── message_repository.go
│           ├── token.go
│           ├── user.go
│           └── user_repository.go
├── pkg
│   ├── encryption.go
│   ├── logger.go
│   ├── passhasher.go
│   └── validator.go
├── README.md
└── templates
    ├── chat.html
    ├── error.html
    ├── login.html
    └── register.html

```

---
# 🚀 Getting Started

## Requirements

Make sure you have installed:

- Go 1.26+
- Docker
- Docker Compose

---

# ⚙️ Configuration

Create a `.env` file in the project root and fill it acording to `.env.example` file:

```env
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=secure_chat_db
JWT_SECRET=your-secret-key
CSRF_SECRET=your-csrf-secret
REDIS_URL=redis://localhost:6379/0
NATS_URL=nats://localhost:4222
CACHE_ENCRYPTION_KEY=your-64-character-hex-key
ADMIN_BOOTSTRAP_USERNAME=admin
ADMIN_BOOTSTRAP_PASSWORD=Admin123@
````

---

# 🐳 Run With Docker

Start the complete application stack:

```bash
docker compose up --build
```

This starts:

* PostgreSQL
* Redis
* NATS JetStream
* Secure Chat API

The application will be available at:

```
http://localhost:8080
```

---

# 💻 Run Locally

Start infrastructure services:

```bash
docker compose up postgres redis nats
```

Run the application:

```bash
go run ./cmd/server
```

---



# 🤝 Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch

```bash
git checkout -b feature/my-feature
```

3. Commit your changes

```bash
git commit -m "Add new feature"
```

4. Push your branch

```bash
git push origin feature/my-feature
```

5. Open a Pull Request

Please make sure your code follows the existing architecture and includes tests when possible.

---

# 📄 License

This project is licensed under the MIT License.

You are free to:

* Use
* Modify
* Distribute
* Commercialize

the software under the conditions of the MIT License.

See the full license text in the [LICENSE](LICENSE) file.

