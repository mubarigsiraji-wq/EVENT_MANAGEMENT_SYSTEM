# =========================
# Stage 1: Build React App
# =========================
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files first for Docker cache
COPY package*.json ./

# Install dependencies
RUN npm install --legacy-peer-deps

# Copy source code
COPY . .

# Build React/Vite application
RUN npm run build


# =========================
# Stage 2: Production Server
# =========================
FROM nginx:alpine AS deployment

# Remove default nginx files
RUN rm -rf /usr/share/nginx/html/*

# Copy React build files
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose HTTP port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
