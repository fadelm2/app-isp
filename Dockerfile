# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o web cmd/web/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler cmd/scheduler/main.go

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binaries and configs from builder
COPY --from=builder /app/web .
COPY --from=builder /app/scheduler .
COPY --from=builder /app/config.json .
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /app/.env .

# Create storage directory for uploads
RUN mkdir -p storage

EXPOSE 9030

CMD ["./web"]
