# Build stage
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM gcr.io/distroless/static-debian11

ENV PORT=8080
EXPOSE 8080

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Run as non-root user
USER nonroot:nonroot

# Command to run
CMD ["./main"]