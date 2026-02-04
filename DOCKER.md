# Docker Deployment Guide

## Quick Start

### 1. Using Docker Compose (Recommended)

The easiest way to run the entire stack:

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes database data)
docker-compose down -v
```

Services will be available at:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432

### 2. First Time Setup

After starting the containers:

1. Access the frontend at http://localhost:3000
2. Click "Register" to create an account
3. The PostgreSQL database starts empty - you'll need to import your GnuCash data

## Importing GnuCash Data

### Option 1: Import from existing PostgreSQL backup

```bash
# Copy your GnuCash backup into the container
docker cp /path/to/gnucash-backup.sql gnucash-db:/tmp/

# Restore the backup
docker exec -i gnucash-db psql -U gnucash -d gnucash < /tmp/gnucash-backup.sql
```

### Option 2: Connect to existing GnuCash database

Edit `docker-compose.yml` and remove the `postgres` service. Update the backend environment variables:

```yaml
  backend:
    environment:
      - DATABASE_HOST=your-existing-host
      - DATABASE_PORT=5432
      - DATABASE_USER=your-user
      - DATABASE_PASSWORD=your-password
      - DATABASE_NAME=gnucash
      - DATABASE_SSLMODE=require
```

## Configuration

### Environment Variables

The backend supports the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_HOST` | localhost | PostgreSQL host |
| `DATABASE_PORT` | 5432 | PostgreSQL port |
| `DATABASE_USER` | gnucash | Database user |
| `DATABASE_PASSWORD` | gnucash | Database password |
| `DATABASE_NAME` | gnucash | Database name |
| `DATABASE_SSLMODE` | disable | SSL mode (disable, require, verify-full) |
| `JWT_SECRET` | (change me) | JWT signing secret (min 32 chars) |
| `SERVER_PORT` | 8080 | Backend server port |

### Customizing docker-compose.yml

**Change ports:**
```yaml
services:
  frontend:
    ports:
      - "80:80"  # Access on port 80 instead of 3000

  backend:
    ports:
      - "8000:8080"  # Backend on port 8000
```

**Persist data with named volumes:**
```yaml
volumes:
  postgres_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /path/to/your/data
```

**Production settings:**
```yaml
services:
  backend:
    environment:
      - DATABASE_SSLMODE=require
      - JWT_SECRET=${JWT_SECRET}  # From .env file
    restart: always

  postgres:
    environment:
      - POSTGRES_PASSWORD=${DB_PASSWORD}
```

Create a `.env` file:
```bash
JWT_SECRET=your-very-secure-random-secret-key-min-32-characters
DB_PASSWORD=your-secure-database-password
```

## Building Images

### Build all services:
```bash
docker-compose build
```

### Build specific service:
```bash
docker-compose build backend
docker-compose build frontend
```

### Build without cache:
```bash
docker-compose build --no-cache
```

## Individual Container Usage

### Backend Only

```bash
# Build
docker build -t gnucash-backend ./backend

# Run
docker run -d \
  --name gnucash-backend \
  -p 8080:8080 \
  -e DATABASE_HOST=host.docker.internal \
  -e DATABASE_USER=gnucash \
  -e DATABASE_PASSWORD=gnucash \
  -e DATABASE_NAME=gnucash \
  -e JWT_SECRET=your-secret-key \
  gnucash-backend
```

### Frontend Only

```bash
# Build with custom API URL
docker build \
  --build-arg REACT_APP_API_BASE_URL=https://api.yourdomain.com/api/v1 \
  -t gnucash-frontend \
  ./frontend

# Run
docker run -d \
  --name gnucash-frontend \
  -p 3000:80 \
  gnucash-frontend
```

## Production Deployment

### Using Docker Compose in Production

1. **Update docker-compose.yml for production:**

```yaml
version: '3.8'

services:
  backend:
    build: ./backend
    environment:
      - DATABASE_HOST=postgres
      - DATABASE_PORT=5432
      - DATABASE_USER=${DB_USER}
      - DATABASE_PASSWORD=${DB_PASSWORD}
      - DATABASE_NAME=gnucash
      - DATABASE_SSLMODE=require
      - JWT_SECRET=${JWT_SECRET}
    restart: always
    depends_on:
      - postgres

  frontend:
    build:
      context: ./frontend
      args:
        - REACT_APP_API_BASE_URL=https://api.yourdomain.com/api/v1
    restart: always

  postgres:
    image: pgvector/pgvector:pg18-alpine
    environment:
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=gnucash
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: always

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - frontend
      - backend
    restart: always

volumes:
  postgres_data:
```

2. **Create production .env file:**

```bash
DB_USER=gnucash_prod
DB_PASSWORD=your-secure-password-here
JWT_SECRET=your-very-secure-jwt-secret-min-32-chars
```

3. **Set proper permissions:**

```bash
chmod 600 .env
```

4. **Deploy:**

```bash
docker-compose -f docker-compose.yml up -d
```

### Health Checks

All services include health checks:

```bash
# Check service health
docker-compose ps

# View health check logs
docker inspect --format='{{json .State.Health}}' gnucash-backend | jq
```

### Monitoring

**View logs:**
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend

# Last 100 lines
docker-compose logs --tail=100
```

**Resource usage:**
```bash
docker stats
```

## Troubleshooting

### Backend can't connect to database

```bash
# Check if postgres is running
docker-compose ps postgres

# Check postgres logs
docker-compose logs postgres

# Test connection
docker exec -it gnucash-backend wget -O- http://localhost:8080/health
```

### Frontend shows connection error

1. Check backend is running: `docker-compose ps backend`
2. Verify API URL in frontend build args
3. Check browser console for CORS errors

### Database initialization fails

```bash
# Reset database
docker-compose down -v
docker-compose up -d

# Or manually recreate
docker-compose stop postgres
docker volume rm gnucash_postgres_data
docker-compose up -d postgres
```

### Permission errors

```bash
# Fix ownership
docker-compose down
sudo chown -R $USER:$USER .
docker-compose up -d
```

### Out of memory

Increase Docker memory limit in Docker Desktop settings, or add to docker-compose.yml:

```yaml
services:
  backend:
    mem_limit: 512m
  frontend:
    mem_limit: 256m
```

## Backup and Restore

### Backup database:
```bash
docker exec gnucash-db pg_dump -U gnucash gnucash > backup-$(date +%Y%m%d).sql
```

### Restore database:
```bash
docker exec -i gnucash-db psql -U gnucash -d gnucash < backup.sql
```

### Backup volumes:
```bash
docker run --rm \
  -v gnucash_postgres_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/postgres-data-$(date +%Y%m%d).tar.gz /data
```

## Updating

### Update images:
```bash
docker-compose pull
docker-compose up -d
```

### Rebuild and update:
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Clean Up

### Remove containers:
```bash
docker-compose down
```

### Remove containers and volumes:
```bash
docker-compose down -v
```

### Remove images:
```bash
docker-compose down --rmi all -v
```

### Clean up Docker system:
```bash
docker system prune -a
```
