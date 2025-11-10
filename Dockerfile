# Dockerfile
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Compile both binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker ./cmd/worker/main.go


# Final minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy compiled binaries
COPY --from=builder /bin/server .
COPY --from=builder /bin/worker .

# Default command to run server (will be overridden in docker-compose)
CMD ["./server"]