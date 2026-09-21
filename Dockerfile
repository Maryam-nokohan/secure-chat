FROM node:22-bookworm-slim AS frontend
WORKDIR /web
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

FROM golang:1.26.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/secure-chat ./cmd/server/main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/secure-chat .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/internal/adapters/secondary/postgres/migrations ./internal/adapters/secondary/postgres/migrations
COPY --from=frontend /web/dist ./frontend/dist

EXPOSE 8080
CMD ["./secure-chat"]