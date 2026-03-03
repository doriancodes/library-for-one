# ---- Build stage ----
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy go mod first (better caching)
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/app .

# Copy static files (important!)
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./app"]
