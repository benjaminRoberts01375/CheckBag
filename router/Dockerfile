# Build frontend with Node.js and npm
FROM node:24-alpine AS builder

WORKDIR /app

# Copy package files from frontend directory
COPY frontend/package.json ./

# Use npm install instead of npm ci (generates lock file)
RUN npm install

# Copy frontend source
COPY frontend/ .

# Build the frontend
RUN npm run build

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
