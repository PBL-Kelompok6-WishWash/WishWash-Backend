# Stage 1: Build the Go Binary
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Stage 2: Minimal runtime image
FROM alpine:3.20

# Install dependencies (timezone support, certificates)
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy uploads folder structure
COPY --from=builder /app/uploads ./uploads

# Expose port 8080
EXPOSE 8080

# Run the app
CMD ["./main"]
