# GnuCash MCP Server - Quick Start

## What is MCP?

The Model Context Protocol (MCP) is an open protocol that enables LLM agents (like Claude, GPT, etc.) to securely interact with your applications and data. This GnuCash MCP Server exposes your financial data through a standardized interface that AI agents can understand and query.

## Quick Start

### Start All Services

```bash
docker compose up -d
```

This starts:
- PostgreSQL database (port 5432)
- HTTP API server (port 8080)
- Frontend (via nginx)
- **MCP Server (port 8081)** ← New!

### Verify MCP Server is Running

```bash
# Check status
docker compose ps mcp-server

# View logs
docker compose logs -f mcp-server
```

You should see:
```
Connected to PostgreSQL successfully
Registered 11 MCP tools
GnuCash MCP Server starting on http://0.0.0.0:8081
```

## Available Tools

The MCP server exposes 11 tools for LLM agents:

**Accounts:**
- `accounts_list` - List all accounts or filter by type
- `accounts_get` - Get account details by GUID
- `accounts_hierarchy` - Get full account tree
- `accounts_balance` - Get account balance

**Transactions:**
- `transactions_list` - List transactions with filters
- `transactions_get` - Get transaction details

**Analytics:**
- `analytics_expenses` - Get expense analysis
- `analytics_income` - Get income analysis
- `analytics_cashflow` - Get cashflow analysis

**Commodities:**
- `commodities_list` - List all currencies
- `commodities_get` - Get commodity details

## Connecting an LLM Agent

### Example: Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "gnucash": {
      "transport": {
        "type": "http",
        "url": "http://localhost:8081"
      }
    }
  }
}
```

Restart Claude Desktop, and you can now ask questions like:
- "What's my current checking account balance?"
- "Show me all expenses from last month"
- "What's my net worth?"

## Architecture

```
LLM Agent (Claude/GPT)
    ↓ HTTP/MCP Protocol
MCP Server :8081
    ↓ PostgreSQL
Database :5432
```

The MCP server:
- Runs independently from the HTTP API server
- Shares the same PostgreSQL database
- Uses the MCP protocol for standardized communication
- Exposes 11 financial data tools

## Configuration

Environment variables (set in `docker-compose.yml`):

```yaml
MCP_PORT: 8081                    # Server port
MCP_SERVER_NAME: gnucash-mcp-server
MCP_SERVER_VERSION: 1.0.0
DATABASE_HOST: postgres
DATABASE_PORT: 5432
DATABASE_USER: gnucash
DATABASE_PASSWORD: gnucash_password
DATABASE_NAME: gnucash
```

## Full Documentation

See `docs/MCP_SERVER.md` for:
- Detailed tool descriptions and parameters
- Security considerations
- Production deployment guide
- Troubleshooting
- Development guide

## What Was Added

### New Files
- `backend/Dockerfile.mcp` - MCP server Docker image
- `backend/cmd/mcp-server/main.go` - MCP server entry point
- `docs/MCP_SERVER.md` - Full documentation
- `README_MCP.md` - This quick start guide

### Modified Files
- `docker-compose.yml` - Added mcp-server service
- `backend/internal/infrastructure/mcp/` - Implemented MCP protocol
- `backend/configs/config.yaml` - Added MCP documentation

### Key Implementation Details
- Uses official Go MCP SDK (`github.com/modelcontextprotocol/go-sdk`)
- HTTP transport with Server-Sent Events (SSE)
- Shares database repositories with HTTP API server
- Stateless architecture - each request is independent
- 11 tools registered and ready for LLM agents

## Stopping the Server

```bash
# Stop just MCP server
docker compose stop mcp-server

# Stop all services
docker compose down
```

## Next Steps

1. ✅ MCP server is running and healthy
2. Configure your LLM agent (Claude Desktop, etc.) to connect
3. Start asking questions about your financial data!
4. Review security settings in `docs/MCP_SERVER.md` for production use

## Support

- Check logs: `docker compose logs mcp-server`
- Health check: `docker compose ps mcp-server`
- Full docs: `docs/MCP_SERVER.md`
