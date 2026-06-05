# syntax=docker/dockerfile:1

FROM node:24-alpine AS frontend-build
WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

FROM golang:1.25-alpine AS backend-build
WORKDIR /app/backend

COPY backend/go.mod ./
RUN go mod download

COPY backend/ ./
COPY --from=frontend-build /app/frontend/dist ./internal/web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/sunapp ./cmd/server

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app
COPY --from=backend-build /app/sunapp ./sunapp

USER app
EXPOSE 8080
ENTRYPOINT ["./sunapp"]
