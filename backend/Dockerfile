# Build backend for production with Go and host with Alpine
FROM golang:alpine AS builder

# Install git
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:3.22

# System group && system user as a part of the system group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /home/appuser

# Copy the binary from builder
COPY --from=builder --chown=appuser:appgroup /app/main /home/appuser/main

# Create the CheckBag directory and set proper permissions
RUN mkdir -p /home/appuser/CheckBag && \
    chown -R appuser:appgroup /home/appuser/CheckBag

# Switch to non-root user
USER appuser

# Run the binary
CMD ["./main"]
