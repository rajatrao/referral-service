# Stage 1: Build the Go binary
FROM golang:1.23.4 AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .
COPY config/base.yaml ./config/


# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Minimal runtime image
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/config/base.yaml ./config/

# Expose the service port (change as needed)
EXPOSE 8090

# Run the service
ENTRYPOINT ["./main"]
