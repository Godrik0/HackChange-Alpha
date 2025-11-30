# Stage 1: Build Angular app
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies with legacy peer deps to resolve Angular 21 conflicts
RUN npm ci --legacy-peer-deps

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Stage 2: Production with Node.js server
FROM node:20-alpine

WORKDIR /app

# Copy built application from builder stage
COPY --from=builder /app/dist/frontend ./dist/frontend

# Copy package files for production dependencies
COPY package*.json ./

# Install only production dependencies with legacy peer deps
RUN npm ci --omit=dev --legacy-peer-deps

# Expose port
EXPOSE 4000

# Start the SSR server
CMD ["node", "dist/frontend/server/server.mjs"]
