# Multi-stage Dockerfile for SatSet Frontend (React + Vite + Tailwind CSS)

# Stage 1: Build production bundle
FROM node:22-alpine AS builder

WORKDIR /app

# Configure npm network timeouts and retry for slow connections
RUN npm config set fetch-retries 5 \
    && npm config set fetch-retry-mintimeout 20000 \
    && npm config set fetch-retry-maxtimeout 300000 \
    && npm config set fetch-timeout 300000

# Install dependencies using clean install
COPY package.json package-lock.json ./
RUN npm ci --no-audit --prefer-offline --maxsockets 5

# Copy source code
COPY . .

# Build arguments for Vite environment variables
ARG VITE_API_BASE_URL=/api
ARG VITE_APP_NAME="SatSet Laundry"
ARG VITE_APP_DESCRIPTION="Sistem Manajemen Laundry & POS"
ARG VITE_APP_TIMEZONE="Asia/Jakarta"
ARG VITE_APP_DEFAULT_LOCALE="id-ID"
ARG VITE_APP_CURRENCY="IDR"

ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
ENV VITE_APP_NAME=$VITE_APP_NAME
ENV VITE_APP_DESCRIPTION=$VITE_APP_DESCRIPTION
ENV VITE_APP_TIMEZONE=$VITE_APP_TIMEZONE
ENV VITE_APP_DEFAULT_LOCALE=$VITE_APP_DEFAULT_LOCALE
ENV VITE_APP_CURRENCY=$VITE_APP_CURRENCY

# Run build
RUN npm run build

# Stage 2: Production web server using Nginx
FROM nginx:alpine

# Copy Nginx server configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy production static assets
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 -O - http://127.0.0.1/ > /dev/null || exit 1

CMD ["nginx", "-g", "daemon off;"]
