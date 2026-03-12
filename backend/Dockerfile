# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy source code and build
COPY . .
RUN go build -o main ./cmd/api/main.go

# Stage 2: Final minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 (Go app default)
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
