# Stage 1: Build frontend với Node + pnpm
FROM node:20-alpine AS builder

WORKDIR /app

# Cài pnpm
RUN npm install -g pnpm

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build arg cho API URL
ARG VITE_API_BASE_URL=http://localhost:8081
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# Build production
RUN pnpm run build

# Stage 2: Serve với Nginx
FROM nginx:alpine

# Copy built files vào nginx html folder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
