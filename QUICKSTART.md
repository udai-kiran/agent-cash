# Quick Start Guide

This guide will get Agentic Cash running in under 5 minutes.

## Prerequisites

- Docker and Docker Compose installed
- Git (to clone the repository)

## Step 1: Clone and Setup

```bash
# Clone the repository
git clone https://github.com/udai-kiran/agentic-cash.git
cd agentic-cash

# Copy the example environment file
cp .env.example .env
```

## Step 2: Configure (Optional)

If you want to customize settings, edit the `.env` file:

```bash
nano .env
```

Key settings to consider:
- `DATABASE_PASSWORD` - Database password
- `JWT_SECRET` - JWT secret key (use a random 32+ character string)

## Step 3: Start Services

```bash
# Start all services in detached mode
docker-compose up -d

# Or use Make
make start
```

This will:
- Start PostgreSQL database
- Start Go backend API
- Start React frontend
- Create necessary database tables

## Step 4: Access the Application

Open your browser and navigate to:

**http://localhost:3000**

## Step 5: Create Your First Account

1. Click "Register" in the top right
2. Enter your email and password (min 8 characters)
3. Click "Create account"
4. You'll be automatically logged in

## Step 6: Import GnuCash Data (Optional)

If you have existing GnuCash data in PostgreSQL:

```bash
# Export from your existing GnuCash database
pg_dump -h your-host -U your-user -d gnucash > gnucash-backup.sql

# Import into Agentic Cash database
docker exec -i gnucash-db psql -U gnucash -d gnucash < gnucash-backup.sql
```

Or if you have an XML file, use GnuCash to export to PostgreSQL first.

## Verify Everything Works

### Check Services

```bash
# View running containers
docker-compose ps

# All three should show "Up" status:
# - gnucash-db
# - gnucash-backend
# - gnucash-frontend
```

### Check Health

```bash
# Backend health check
curl http://localhost:8080/health

# Should return: {"status":"ok"}
```

### Test in Browser

1. **Dashboard** - Should show Net Worth and Income/Expense charts
2. **Accounts** - Click to see your account hierarchy
3. **Transactions** - View and filter transactions
4. **Analytics** - See detailed financial analytics

## Common Issues

### Port Already in Use

If you get "port already in use" errors:

```bash
# Change ports in docker-compose.yml
services:
  frontend:
    ports:
      - "8000:80"  # Change 3000 to 8000
  backend:
    ports:
      - "8081:8080"  # Change 8080 to 8081
```

### Database Connection Failed

```bash
# Check database logs
docker-compose logs postgres

# Restart the backend
docker-compose restart backend
```

### No Data Showing

If the dashboard is empty:
- You need to import GnuCash data (see Step 6)
- Or add test data manually through GnuCash application

## Next Steps

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
```

### Stop Services

```bash
# Stop but keep data
docker-compose down

# Stop and remove all data (WARNING: deletes database)
docker-compose down -v
```

### Backup Database

```bash
# Create backup
make backup-db

# Or manually
docker exec gnucash-db pg_dump -U gnucash gnucash > backup.sql
```

### Development Mode

If you want to modify the code:

```bash
# Start in development mode with hot reload
make dev

# Or
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## Useful Commands

```bash
# Start services
make start

# View logs
make logs

# Stop services
make stop

# Check health
make health

# Backup database
make backup-db

# Access database shell
make shell-db

# Access backend shell
make shell-backend

# Clean everything
make clean
```

## Getting Help

- **Full documentation**: See [README.md](README.md)
- **Docker guide**: See [DOCKER.md](DOCKER.md)
- **Testing guide**: See [TESTING.md](TESTING.md)
- **Issues**: https://github.com/udai-kiran/agentic-cash/issues

## Success!

If you can:
1. Access http://localhost:3000
2. Register and login
3. See the dashboard

You're all set! Enjoy using Agentic Cash for your GnuCash data visualization.
