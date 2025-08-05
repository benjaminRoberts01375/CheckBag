# Build frontend with bun/vite
FROM oven/bun:1.2.19 AS builder

WORKDIR /app

# Copy package files from frontend directory
COPY frontend/package.json frontend/bun.lock ./

# Install dependencies
RUN bun install

# Copy frontend source
COPY frontend/ .

# Build the frontend
RUN bun run build

# Setup NGINX for production
FROM nginx:1.27.4-alpine-slim

# Remove the default nginx configuration
RUN rm /etc/nginx/conf.d/default.conf

# Copy built frontend files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx configuration
COPY router/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
