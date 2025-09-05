# Build frontend
FROM oven/bun:latest AS frontend-builder

WORKDIR /app

# Copy package files from frontend directory
COPY ./frontend/package.json ./frontend/bun.lock ./

# Install dependencies
RUN bun install

# Copy frontend source
COPY ./frontend/ .

# Build the frontend
RUN bun run build

# Build backend for production
FROM golang:alpine AS backend-builder

# Install git
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY ./backend/go.mod ./backend/go.sum ./

# Download dependencies
RUN go mod download && go mod tidy

# Copy backend source code
COPY ./backend/ .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:3.22

# System group && system user as a part of the system group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /home/appuser

# Copy the Go binary from backend builder
COPY --from=backend-builder --chown=appuser:appgroup /app/main /home/appuser/main

# Copy the frontend static files from frontend builder
COPY --from=frontend-builder --chown=appuser:appgroup /app/dist /home/appuser/static

# Create the CheckBag directory and make it world-writable to handle volume mounts
RUN mkdir -p /home/appuser/CheckBag && \
    chmod 777 /home/appuser/CheckBag

# Switch to non-root user
USER appuser

# Run the binary
CMD ["./main"]
