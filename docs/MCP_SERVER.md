# GnuCash MCP Server Documentation

## Overview

The GnuCash MCP (Model Context Protocol) Server provides LLM agents with structured access to your GnuCash financial data. It exposes account information, transactions, analytics, and commodity data through the standardized MCP protocol.

## Architecture

```
┌─────────────────┐
│   LLM Agent     │
│  (Claude, GPT)  │
└────────┬────────┘
         │ HTTP/SSE
         │ Port 8081
         ▼
┌─────────────────────────────────┐
│     MCP Server Container        │
│  ┌───────────────────────────┐  │
│  │  MCP Protocol Handler     │  │
│  ├───────────────────────────┤  │
│  │  11 Available Tools:      │  │
│  │  • accounts_list          │  │
│  │  • accounts_get           │  │
│  │  • accounts_hierarchy     │  │
│  │  • accounts_balance       │  │
│  │  • transactions_list      │  │
│  │  • transactions_get       │  │
│  │  • analytics_expenses     │  │
│  │  • analytics_income       │  │
│  │  • analytics_cashflow     │  │
│  │  • commodities_list       │  │
│  │  • commodities_get        │  │
│  └───────────────────────────┘  │
└─────────────┬───────────────────┘
              │ PostgreSQL
              ▼
┌─────────────────────────────────┐
│     PostgreSQL Container        │
└─────────────────────────────────┘
```

## Starting the MCP Server

### Using Docker Compose (Recommended)

```bash
# Start all services including MCP server
docker-compose up -d

# Start only the MCP server (requires postgres to be running)
docker-compose up -d mcp-server

# View MCP server logs
docker-compose logs -f mcp-server
```

### Standalone Build

```bash
cd backend
go build ./cmd/mcp-server
./mcp-server
```

## Configuration

The MCP server uses environment variables for configuration:

| Variable | Default | Description |
|----------|---------|-------------|
| `MCP_PORT` | `8081` | Port for MCP server to listen on |
| `MCP_SERVER_NAME` | `gnucash-mcp-server` | Server name in MCP protocol |
| `MCP_SERVER_VERSION` | `1.0.0` | Server version in MCP protocol |
| `DATABASE_HOST` | `postgres` | PostgreSQL host |
| `DATABASE_PORT` | `5432` | PostgreSQL port |
| `DATABASE_USER` | `gnucash` | Database user |
| `DATABASE_PASSWORD` | `gnucash_password` | Database password |
| `DATABASE_NAME` | `gnucash` | Database name |
| `DATABASE_SSLMODE` | `disable` | SSL mode for database connection |
| `GO_ENV` | - | Set to `production` for production logging |

## Available Tools

### Account Tools

#### `accounts_list`
Lists all accounts or filters by account type.

**Parameters:**
- `type` (optional): Account type filter (ASSET, LIABILITY, EQUITY, INCOME, EXPENSE, etc.)

**Example Response:**
```json
{
  "accounts": [
    {
      "guid": "abc-123",
      "name": "Checking Account",
      "type": "ASSET",
      "balance": "1500.00",
      "commodity": "USD"
    }
  ],
  "count": 1
}
```

#### `accounts_get`
Gets detailed information about a specific account.

**Parameters:**
- `guid` (required): Account GUID

#### `accounts_hierarchy`
Returns the complete account hierarchy as a tree structure.

**Parameters:** None

#### `accounts_balance`
Gets the current balance of a specific account.

**Parameters:**
- `guid` (required): Account GUID

### Transaction Tools

#### `transactions_list`
Lists transactions with optional filters.

**Parameters:**
- `account_guid` (optional): Filter by account GUID
- `start_date` (optional): Start date in YYYY-MM-DD format
- `end_date` (optional): End date in YYYY-MM-DD format
- `description` (optional): Filter by description (partial match)

#### `transactions_get`
Gets detailed information about a specific transaction.

**Parameters:**
- `guid` (required): Transaction GUID

### Analytics Tools

#### `analytics_expenses`
Get expense analysis for a date range.

**Parameters:**
- `start_date` (optional): Start date in YYYY-MM-DD format (defaults to 1 month ago)
- `end_date` (optional): End date in YYYY-MM-DD format (defaults to today)

#### `analytics_income`
Get income analysis for a date range.

**Parameters:**
- `start_date` (optional): Start date in YYYY-MM-DD format
- `end_date` (optional): End date in YYYY-MM-DD format

#### `analytics_cashflow`
Get cash flow analysis showing income and expenses over time.

**Parameters:**
- `start_date` (optional): Start date in YYYY-MM-DD format
- `end_date` (optional): End date in YYYY-MM-DD format

### Commodity Tools

#### `commodities_list`
Lists all commodities (currencies) in the database.

**Parameters:** None

#### `commodities_get`
Gets detailed information about a specific commodity.

**Parameters:**
- `guid` (required): Commodity GUID

## Connecting LLM Agents

### Claude Desktop

Add to your Claude Desktop configuration (`claude_desktop_config.json`):

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

### Generic MCP Client (Python)

```python
from modelcontextprotocol import Client
from modelcontextprotocol.transports.http import HTTPTransport

# Create HTTP transport
transport = HTTPTransport("http://localhost:8081")

# Create and connect client
client = Client(transport)
await client.connect()

# List available tools
tools = await client.list_tools()
print("Available tools:", [t.name for t in tools])

# Call a tool
result = await client.call_tool("accounts_list", {})
print(result)
```

### Generic MCP Client (TypeScript/JavaScript)

```typescript
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { StreamableHTTPClientTransport } from '@modelcontextprotocol/sdk/client/http.js';

// Create HTTP transport
const transport = new StreamableHTTPClientTransport({
  endpoint: 'http://localhost:8081'
});

// Create and connect client
const client = new Client({
  name: 'gnucash-client',
  version: '1.0.0'
}, {});

await client.connect(transport);

// List available tools
const tools = await client.listTools();
console.log('Available tools:', tools.tools.map(t => t.name));

// Call a tool
const result = await client.callTool({
  name: 'accounts_list',
  arguments: {}
});
console.log(result);
```

## Security Considerations

### Current Implementation
- **No authentication**: The MCP server currently has no built-in authentication
- **Network access**: Exposed on port 8081 by default

### Recommendations for Production

1. **Network Isolation**: Run MCP server on internal network only
   ```yaml
   # docker-compose.yml
   mcp-server:
     ports:
       # Remove this to keep it internal-only
       # - "8081:8081"
   ```

2. **Reverse Proxy with Authentication**: Use nginx or traefik with authentication
   ```nginx
   location /mcp {
       auth_basic "MCP Server";
       auth_basic_user_file /etc/nginx/.htpasswd;
       proxy_pass http://mcp-server:8081;
   }
   ```

3. **VPN Access**: Require VPN connection to access MCP server

4. **API Key Authentication**: Implement in future version

## Troubleshooting

### MCP Server Won't Start

Check logs:
```bash
docker-compose logs mcp-server
```

Common issues:
- Database not ready: Wait for postgres health check to pass
- Port already in use: Change `MCP_PORT` environment variable
- Database connection error: Verify `DATABASE_*` environment variables

### Tools Not Discoverable

Verify server is responding:
```bash
curl http://localhost:8081
```

### Connection Refused from LLM Agent

- Ensure MCP server is running: `docker-compose ps mcp-server`
- Check firewall rules allow connections to port 8081
- Verify network configuration if running on remote server

## Health Monitoring

The MCP server includes a health check that verifies the server is listening on the configured port.

Check health status:
```bash
# Using Docker
docker inspect --format='{{.State.Health.Status}}' gnucash-mcp-server

# Using docker-compose
docker-compose ps mcp-server
```

## Development

### Running Tests

```bash
cd backend
go test ./internal/infrastructure/mcp/...
```

### Adding New Tools

1. Define tool parameters struct in appropriate `*_tools.go` file
2. Implement handler function with signature:
   ```go
   func (s *MCPServer) handleToolName(
       ctx context.Context,
       req *mcp.CallToolRequest,
       params *ToolParams,
   ) (*mcp.CallToolResult, any, error)
   ```
3. Register tool in `server.go` `registerTools()` method
4. Build and test

### Building Locally

```bash
cd backend
go build ./cmd/mcp-server
./mcp-server
```

## API Endpoint

The MCP server exposes a single HTTP endpoint that handles all MCP protocol messages:

- **URL**: `http://localhost:8081/mcp/message`
- **Method**: POST
- **Content-Type**: application/json
- **Protocol**: MCP over HTTP with Server-Sent Events (SSE)

## License

This MCP server implementation is part of the GnuCash project and follows the same license.

## Support

For issues or questions:
- Check logs: `docker-compose logs mcp-server`
- Review configuration in `docker-compose.yml`
- Verify environment variables are set correctly
