# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build tools
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o engagesync .
RUN CGO_ENABLED=0 GOOS=linux go build -o migration ./db/migration.go

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata postgresql-client

COPY --from=builder /app/engagesync .
COPY --from=builder /app/migration .
COPY --from=builder /app/db/migrations ./db/migrations
COPY start.sh .
COPY .env .

RUN chmod +x start.sh

EXPOSE 8000

CMD ["./start.sh"]
