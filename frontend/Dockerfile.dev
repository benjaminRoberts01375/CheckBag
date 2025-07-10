FROM oven/bun:latest

WORKDIR /app

# Copy package files
COPY package.json bun.lockb* ./

# Install dependencies
RUN bun install

# Copy project files
COPY . .

# Expose port
EXPOSE 5173

# Start development server with host binding to allow external access
CMD ["bun", "run", "dev"]
