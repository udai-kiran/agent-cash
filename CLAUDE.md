# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Agentic Cash is a financial dashboard and AI chat interface for GnuCash data. It connects to an existing GnuCash PostgreSQL database and provides a web UI, REST API, and conversational AI for querying financial data.

## Services & Ports

| Service | Language | Port | Location |
|---------|----------|------|----------|
| Frontend | React 19 + TypeScript | 3000 | `/frontend` |
| Backend API | Go (Gin) | 8080 | `/backend` |
| MCP Server | Go | 8081 | `/backend/cmd/mcp-server` |
| Chainlit Chat UI | Python | 8082 | `/chainlit` |
| Agents Service | Python (FastAPI) | 8083 | `/agents` |
| PostgreSQL | - | 5432 | Docker |

## Common Commands

### Docker (primary workflow)
```bash
make start          # Start all services
make stop           # Stop all services
make logs           # View all logs
make logs-backend   # Backend logs only
make health         # Check service health
make ps             # Show running containers
make shell-db       # Open psql shell
make backup-db      # Backup database
```

### Backend (Go)
```bash
cd backend
go build -o bin/server ./cmd/server    # Build API server
go build -o bin/mcp ./cmd/mcp-server   # Build MCP server
go test ./...                          # Run all tests
go test ./internal/application/...     # Run specific package tests
```

### Frontend (React)
```bash
cd frontend
npm install         # Install dependencies
npm start           # Dev server on :3000
npm run build       # Production build
npm test            # Run tests
```

### Agents (Python)
```bash
cd agents
uv sync             # Install dependencies
uv run uvicorn main:app --reload --port 8083
```

## Architecture

### Backend: Clean Architecture (DDD)

```
backend/internal/
├── domain/           # Entities & repository interfaces (no external deps)
│   ├── entity/       # Account, Transaction, User, Split, Commodity
│   └── repository/   # Interface definitions
├── application/      # Use cases, services, DTOs
│   ├── service/      # AuthService, AccountService, AnalyticsService
│   └── dto/          # Request/response objects
├── infrastructure/   # Implementations
│   ├── persistence/postgres/  # Repository implementations (pgx, no ORM)
│   ├── auth/         # JWT token manager
│   └── mcp/          # MCP server with 11 financial tools
└── interfaces/http/  # Gin handlers & middleware
    ├── handler/      # HTTP request handlers
    └── middleware/    # Auth, CORS, rate limiting, logging
```

Dependency injection is manual in `backend/cmd/server/main.go`. Database access uses pgx directly (no ORM). Financial calculations use `shopspring/decimal` for precision.

The backend creates app-specific tables (`users`, `refresh_tokens`, `app_config`) on startup via `postgres.InitializeAppTables()`, while GnuCash tables (`accounts`, `transactions`, `splits`, `commodities`, `prices`) are read-only from the existing database.

### Frontend: React + TanStack Query

- **Server state**: TanStack Query (React Query) via custom hooks in `src/hooks/`
- **Auth state**: React Context (`src/context/AuthContext.tsx`) with localStorage tokens
- **API layer**: Axios with interceptors for auto token refresh (`src/api/client.ts`)
- **Routing**: React Router v7
- **Styling**: Tailwind CSS
- **Charts**: Recharts

### AI Chat Pipeline

```
Chainlit UI (8082) → Agents Service (8083) → MCP Server (8081) → PostgreSQL
```

The agents service uses Strands framework with an LLM (configurable, default GPT-4) and connects to the MCP server which exposes 11 GnuCash data tools (account listing, balance queries, transaction search, analytics).

## API Routes

All REST endpoints are prefixed with `/api/v1/`. Auth endpoints (`/auth/*`) are rate-limited to 5 req/min. Protected routes use JWT Bearer tokens (15min access, 7d refresh).

Key route groups: `/auth` (register/login/refresh/logout), `/accounts` (list/hierarchy/balance), `/transactions` (list/get with filtering), `/commodities`, `/analytics` (income-expense/category-breakdown/net-worth).

Router definition: `backend/internal/interfaces/http/router.go`

## Environment

Copy `.env.example` to `.env`. Key variables:
- `DATABASE_*` - PostgreSQL connection
- `JWT_SECRET` - Min 32 chars for token signing
- `OPENAI_API_KEY` - Required for chat/agents service
- `REACT_APP_API_BASE_URL` - Frontend API target (default `http://localhost:8080/api/v1`)

Backend also reads `backend/configs/config.yaml` (Viper-based, env vars override).

## Python Conventions

Per `.cursor/agents/python.md`: target Python 3.14, use modern typing (`list[str]`, `|` unions), prefer `pytest`, `uv` for deps, `pydantic-settings` for config, `asyncio` for I/O concurrency. Validate data at service boundaries with Pydantic models.
