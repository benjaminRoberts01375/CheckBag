# Build backend for production with Go and host with Alpine
FROM golang:alpine AS builder

# Install git
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Mount GitHub token at build time
RUN --mount=type=secret,id=git_env \
    if [ -f "/run/secrets/git_env" ]; then \
    TOKEN=$(cat /run/secrets/git_env) && \
    git config --global url."https://${TOKEN}@github.com/".insteadOf "https://github.com/"; \
    fi

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

# Create configs directory and set permissions
USER appuser
# Run the binary
CMD ["./main"]
