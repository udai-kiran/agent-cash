# Agentic Cash

A modern financial dashboard for GnuCash built with Go and React.

## Overview

Agentic Cash provides a web-based interface to visualize and analyze your GnuCash financial data. It connects to your existing GnuCash PostgreSQL database and provides:

- Interactive account hierarchy visualization
- Transaction browsing and filtering
- Financial analytics and charts
- Budget tracking and reporting

## Architecture

The project consists of two main components:

### Backend (Go)
- Clean architecture with domain-driven design
- PostgreSQL connection to GnuCash database
- RESTful API with JWT authentication
- Located in `/backend`

### Frontend (React + TypeScript)
- Modern React with TypeScript
- TanStack Query for server state
- Tailwind CSS for styling
- Recharts for visualizations
- Located in `/frontend`

## Quick Start

### Option 1: Docker (Recommended)

**Prerequisites:** Docker and Docker Compose

```bash
# Clone the repository
git clone https://github.com/udai-kiran/agentic-cash.git
cd agentic-cash

# Copy environment file
cp .env.example .env

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

**Or using Make:**
```bash
make start      # Start services
make logs       # View logs
make stop       # Stop services
```

Services will be available at:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Chat Interface**: http://localhost:8082 (see [Chat Setup](#chat-interface) below)
- **PostgreSQL**: localhost:5432

See [DOCKER.md](DOCKER.md) for detailed Docker documentation.

### Option 2: Local Development

**Prerequisites:**
- Go 1.25+
- Node.js 25.5.0+
- PostgreSQL with GnuCash database

**Backend:**
```bash
cd backend
go build -o bin/server ./cmd/server
# Edit configs/config.yaml with your database credentials
./bin/server
```

**Frontend:**
```bash
cd frontend
npm install
# Create .env file with REACT_APP_API_BASE_URL=http://localhost:8080/api/v1
npm start
```

Backend runs on http://localhost:8080
Frontend runs on http://localhost:3000

## Project Structure

```
gnucash/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── domain/          # Business entities and interfaces
│   │   ├── application/     # Use cases and DTOs
│   │   ├── infrastructure/  # Database and external services
│   │   └── interfaces/      # HTTP handlers and middleware
│   ├── pkg/                 # Shared utilities
│   └── configs/             # Configuration files
└── frontend/
    └── src/
        ├── api/             # API client
        ├── components/      # React components
        ├── hooks/           # Custom hooks
        ├── pages/           # Page components
        ├── types/           # TypeScript types
        └── utils/           # Utility functions
```

## Features

### Completed

- [x] **Phase 1: Foundation**
  - [x] Clean architecture backend with Go
  - [x] Account hierarchy API
  - [x] React frontend with TypeScript
  - [x] Account tree visualization with collapsible nodes

- [x] **Phase 2: Authentication**
  - [x] JWT authentication system
  - [x] User registration and login
  - [x] Token refresh mechanism
  - [x] Protected routes (optional)

- [x] **Phase 3: Transactions**
  - [x] Transaction API with filtering
  - [x] Transaction table with date/description filters
  - [x] Pagination support
  - [x] Split viewing

- [x] **Phase 4: Analytics**
  - [x] Income vs expense charts
  - [x] Category breakdown (pie charts)
  - [x] Net worth calculation
  - [x] Interactive dashboard

### Roadmap

- [ ] **Phase 5: Budgets**
  - [ ] Budget tracking
  - [ ] Budget vs actual reports
  - [ ] Budget progress visualization

- [ ] **Future Enhancements**
  - [ ] Transaction creation/editing
  - [ ] CSV export
  - [ ] Multi-currency support
  - [ ] Scheduled reports

## Development Status

The application is feature-complete for phases 1-4. It provides:
- Full account hierarchy browsing
- Transaction viewing with filters
- User authentication
- Comprehensive financial analytics
- Interactive charts and visualizations

## Chat Interface

A conversational AI interface for querying your financial data using natural language!

**Quick Start:**
```bash
# Add your OpenAI API key
cp .env.example .env
# Edit .env and add: OPENAI_API_KEY=sk-your-key-here

# Start all services (including chat)
docker-compose up -d

# Open chat interface
open http://localhost:8082
```

**Example queries:**
- "What's my checking account balance?"
- "Show me expenses from last month"
- "What's my net worth?"
- "Find all transactions over $500"

**Documentation:**
- **[QUICKSTART_CHAT.md](QUICKSTART_CHAT.md)** - 5-minute setup guide
- **[README_CHAT.md](README_CHAT.md)** - Complete chat documentation
- **[CHAT_REFERENCE.md](CHAT_REFERENCE.md)** - Quick reference

## Documentation

- **[DOCKER.md](DOCKER.md)** - Docker setup and deployment guide
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Production deployment instructions
- **[TESTING.md](TESTING.md)** - Testing and validation guide
- **[backend/README.md](backend/README.md)** - Backend API documentation
- **[frontend/README.md](frontend/README.md)** - Frontend development guide
- **[README_CHAT.md](README_CHAT.md)** - Chat interface guide
- **[README_MCP.md](README_MCP.md)** - MCP server documentation

## Quick Commands

```bash
# Using Make
make start          # Start all services
make stop           # Stop all services
make logs           # View logs
make dev            # Start development environment
make backup-db      # Backup database
make health         # Check service health

# Using Docker Compose
docker-compose up -d              # Start services
docker-compose down               # Stop services
docker-compose logs -f            # View logs
docker-compose ps                 # List containers
```

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.
# agent-cash
