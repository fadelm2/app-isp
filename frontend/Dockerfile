# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package management files
COPY package*.json ./
RUN npm install

# Copy source code and build
COPY . .
RUN npm run build

# Serve stage
FROM nginx:alpine

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
