# Deployment Guide

## Prerequisites

1. PostgreSQL database with GnuCash data
2. Go 1.21+ installed
3. Node.js 25.5.0+ installed

## Backend Deployment

### 1. Configure Database

Edit `backend/configs/config.yaml`:

```yaml
server:
  port: 8080
  readTimeout: 15s
  writeTimeout: 15s

database:
  host: your-db-host
  port: 5432
  user: your-db-user
  password: your-db-password
  dbname: gnucash
  sslmode: require  # Use 'require' in production
  maxConns: 10
  minConns: 2

jwt:
  secret: your-secure-secret-key-min-32-chars
  accessTokenTTL: 15m
  refreshTokenTTL: 168h
```

### 2. Build Backend

```bash
cd backend
go build -o bin/server ./cmd/server
```

### 3. Run Backend

```bash
./bin/server
```

The server will:
- Start on port 8080
- Create `app_users` and `refresh_tokens` tables automatically
- Connect to your GnuCash PostgreSQL database (read-only for GnuCash tables)

## Frontend Deployment

### 1. Configure API URL

Create `frontend/.env.production`:

```
REACT_APP_API_BASE_URL=https://your-api-domain.com/api/v1
```

### 2. Build Frontend

```bash
cd frontend
npm install
npm run build
```

This creates a production build in `frontend/build/`.

### 3. Serve Frontend

Option A: Use a static file server (nginx, Apache)
Option B: Deploy to Vercel, Netlify, or similar

Example nginx configuration:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /path/to/frontend/build;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## Docker Deployment (Optional)

### Backend Dockerfile

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs
CMD ["./server"]
```

### Frontend Dockerfile

```dockerfile
FROM node:25.5.0-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

### Docker Compose

```yaml
version: '3.8'

services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_HOST=postgres
      - DATABASE_PORT=5432
      - DATABASE_USER=gnucash
      - DATABASE_PASSWORD=secret
      - DATABASE_NAME=gnucash
    depends_on:
      - postgres

  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend

  postgres:
    image: pgvector/pgvector:pg18-alpine
    environment:
      - POSTGRES_USER=gnucash
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=gnucash
    volumes:
      - postgres-data:/var/lib/postgresql/data

volumes:
  postgres-data:
```

## Security Considerations

1. **JWT Secret**: Use a strong, random secret (min 32 characters)
2. **HTTPS**: Always use HTTPS in production
3. **Database**: Use SSL for database connections
4. **CORS**: Configure CORS to only allow your frontend domain
5. **Rate Limiting**: Consider adding rate limiting middleware
6. **Environment Variables**: Never commit secrets to version control

## Monitoring

Consider adding:
- Application logs (structured logging)
- Health check endpoints (already available at `/health`)
- Metrics (Prometheus)
- Error tracking (Sentry)

## Backup

Regularly backup:
1. GnuCash PostgreSQL database
2. Application user database (`app_users`, `refresh_tokens`)
3. Configuration files
