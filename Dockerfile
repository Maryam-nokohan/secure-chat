FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copies everything (including internal/ and your local .env file)
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/secure-chat ./cmd/server/main.go

FROM alpine:3.22
WORKDIR /app

# Pull the binary, env, templates, AND migrations into the final runner stage
COPY --from=builder /app/secure-chat .
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/templates ./templates

# CRITICAL: Copy the migrations folder over so the app can find schema.sql
COPY --from=builder /app/internal/adapters/secondary/postgres/migrations ./internal/adapters/secondary/postgres/migrations

EXPOSE 8080

CMD ["./secure-chat"]