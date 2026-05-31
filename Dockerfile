# Stage 1: Build the Go application
FROM golang:1.26 AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary statically linked, stripped of debug info for size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o brisk ./cmd/brisk/main.go

# Stage 2: Final lightweight image
FROM scratch
COPY --from=builder /app/brisk /brisk

# Set the binary as the entrypoint
ENTRYPOINT ["/brisk"]