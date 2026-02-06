# Chat Interface Architecture

## System Overview

The GnuCash chat interface consists of three main services that work together to provide a conversational AI interface for financial data.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        User's Browser                            │
│                    http://localhost:8082                         │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            │ WebSocket + HTTP
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Chainlit Service (Port 8082)                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Chainlit Framework                                       │  │
│  │  - FastAPI Backend                                        │  │
│  │  - React Frontend                                         │  │
│  │  - WebSocket Handler                                      │  │
│  │  - Session Management (UUID per user)                    │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Agents Client                                            │  │
│  │  - HTTP client to agents service                          │  │
│  │  - Error handling & retry logic                           │  │
│  └───────────────────────────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            │ HTTP POST /api/v1/chat
                            │ { message, session_id }
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Agents Service (Port 8083)                    │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  FastAPI Server                                           │  │
│  │  - REST API endpoints                                     │  │
│  │  - Request validation (Pydantic)                          │  │
│  │  - CORS middleware                                        │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Agent Service                                            │  │
│  │  - Session-based agent management                         │  │
│  │  - One agent instance per user session                    │  │
│  │  - Agent lifecycle management                             │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Finance Agent (Strands)                                  │  │
│  │  - Powered by GPT-4 (or configured model)                 │  │
│  │  - System prompt for financial context                    │  │
│  │  - Tool selection & execution                             │  │
│  │  - Natural language generation                            │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  MCP Client                                               │  │
│  │  - Model Context Protocol client                          │  │
│  │  - Connects to MCP server                                 │  │
│  │  - Tool discovery & invocation                            │  │
│  └───────────────────────────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            │ MCP Protocol (HTTP)
                            │ Tool calls: get_accounts, get_balance, etc.
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    MCP Server (Port 8081)                        │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Go HTTP Server                                           │  │
│  │  - MCP protocol handler                                   │  │
│  │  - Tool registration & discovery                          │  │
│  │  - Request routing                                        │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  11 Financial Tools                                       │  │
│  │  1. get_accounts                                          │  │
│  │  2. get_account_by_id                                     │  │
│  │  3. get_account_balance                                   │  │
│  │  4. get_transactions                                      │  │
│  │  5. get_transaction_by_id                                 │  │
│  │  6. calculate_income_expenses                             │  │
│  │  7. get_account_hierarchy                                 │  │
│  │  8. search_accounts                                       │  │
│  │  9. search_transactions                                   │  │
│  │ 10. get_commodity_prices                                  │  │
│  │ 11. get_budget_info (if budgets exist)                    │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Database Client (GORM)                                   │  │
│  │  - PostgreSQL connection pool                             │  │
│  │  - Query builder                                          │  │
│  │  - Transaction management                                 │  │
│  └───────────────────────────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            │ SQL Queries
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                  PostgreSQL Database (Port 5432)                 │
│                                                                  │
│  Tables:                                                         │
│  - accounts (financial accounts)                                 │
│  - transactions (financial transactions)                         │
│  - splits (transaction line items)                               │
│  - commodities (currencies)                                      │
│  - prices (exchange rates)                                       │
│  - budgets (budget data)                                         │
└─────────────────────────────────────────────────────────────────┘
```

## Request Flow Example

### User asks: "What's my checking account balance?"

```
1. Browser → Chainlit (WebSocket)
   User types message in chat interface

2. Chainlit → Agents Service (HTTP POST)
   POST /api/v1/chat
   {
     "message": "What's my checking account balance?",
     "session_id": "550e8400-e29b-41d4-a716-446655440000"
   }

3. Agents Service → Finance Agent
   AgentService.process_message()
   - Retrieves or creates Finance Agent for session
   - Passes message to agent

4. Finance Agent → OpenAI API
   - Sends prompt with user message
   - Includes system prompt and available tools
   - GPT-4 analyzes query

5. GPT-4 → Finance Agent
   - Decides to call tool: "get_account_balance"
   - Provides parameters: { "account_name": "Checking" }

6. Finance Agent → MCP Client → MCP Server
   Tool call: get_account_balance(account_name="Checking")

7. MCP Server → PostgreSQL
   SELECT balance FROM accounts WHERE name LIKE '%Checking%'

8. PostgreSQL → MCP Server
   Returns: { "balance": 5234.56, "currency": "USD" }

9. MCP Server → Finance Agent
   Tool result: Account balance data

10. Finance Agent → OpenAI API
    Sends tool result to generate natural language response

11. GPT-4 → Finance Agent
    Generated response: "Your checking account balance is $5,234.56"

12. Finance Agent → Agents Service
    Returns formatted response

13. Agents Service → Chainlit (HTTP Response)
    {
      "response": "Your checking account balance is $5,234.56",
      "session_id": "550e8400-e29b-41d4-a716-446655440000"
    }

14. Chainlit → Browser (WebSocket)
    Displays message in chat interface
```

**Total time**: Typically 5-30 seconds depending on:
- Model used (GPT-4 vs GPT-3.5-turbo)
- Query complexity
- Number of tool calls
- Database query time

## Component Details

### Chainlit Service

**Technology**: Python 3.11, Chainlit 2.9.6+, FastAPI
**Purpose**: User-facing chat interface
**Key Features**:
- WebSocket for real-time communication
- Session management with UUIDs
- Error handling and user-friendly messages
- Health check for agents service

**Files**:
- `app.py` - Main Chainlit application
- `utils/agents_client.py` - HTTP client for agents service
- `.chainlit/config.toml` - Chainlit configuration
- `Dockerfile` - Container definition

### Agents Service

**Technology**: Python 3.11, FastAPI, Strands Agents, MCP SDK
**Purpose**: AI agent orchestration and MCP integration
**Key Features**:
- Session-based agent instances
- Strands framework for agent management
- MCP client for tool access
- RESTful API endpoints

**Files**:
- `main.py` - FastAPI application
- `agents/finance_agent.py` - Main financial assistant agent
- `services/agent_service.py` - Agent lifecycle management
- `services/mcp_client.py` - MCP server connection
- `routes/chat.py` - API endpoints
- `config.py` - Configuration management
- `Dockerfile` - Container definition

### MCP Server

**Technology**: Go 1.23, Gorilla Mux, GORM
**Purpose**: Financial data access layer
**Key Features**:
- 11 specialized financial tools
- Model Context Protocol implementation
- Direct database access
- Efficient query optimization

**Location**: `backend/cmd/mcp-server/` and `backend/internal/infrastructure/mcp/`

## Data Flow Patterns

### Simple Query (Single Tool Call)

```
User: "List all accounts"
  → Agent calls: get_accounts()
  → Returns: List of accounts
  → Agent formats: "Here are your accounts: ..."
```

### Complex Query (Multiple Tool Calls)

```
User: "What's my net worth?"
  → Agent calls: get_accounts()
  → Agent filters: Assets and Liabilities
  → Agent calls: get_account_balance() for each
  → Agent calculates: Total Assets - Total Liabilities
  → Agent formats: "Your net worth is $..."
```

### Analytical Query (Aggregation)

```
User: "Show my top 5 expenses"
  → Agent calls: get_transactions(account_type="Expense")
  → Agent calls: calculate_income_expenses()
  → Agent analyzes: Groups and sorts by amount
  → Agent formats: Ranked list with amounts
```

## Session Management

Each user gets a unique session:

```python
# Chainlit creates session
session_id = uuid.uuid4()  # e.g., "550e8400-..."

# Session stored in Chainlit
cl.user_session.set("session_id", session_id)

# Agents service maintains agent per session
agents[session_id] = FinanceAgent()
```

Benefits:
- Conversation context preserved
- Isolated agent instances
- No cross-user data leakage
- Easy cleanup on session end

## Security Considerations

1. **API Keys**: Stored in environment variables, never in code
2. **Network Isolation**: Services communicate via internal Docker network
3. **No Authentication**: Currently no user auth (future enhancement)
4. **CORS**: Allows all origins (should be restricted in production)
5. **Input Validation**: Pydantic models validate all inputs
6. **Error Handling**: Errors sanitized before showing to users

## Scaling Considerations

### Current Architecture (Single Host)

- All services on one machine
- Suitable for: Personal use, small teams (<10 users)
- Resource usage: ~2-4GB RAM, 2 CPU cores

### Future Scaling Options

**Horizontal Scaling**:
```
Load Balancer
    ├─ Chainlit Instance 1
    ├─ Chainlit Instance 2
    └─ Chainlit Instance 3
         ↓
Load Balancer
    ├─ Agents Instance 1
    ├─ Agents Instance 2
    └─ Agents Instance 3
         ↓
MCP Server Cluster
    ├─ MCP Instance 1
    ├─ MCP Instance 2
    └─ MCP Instance 3
         ↓
PostgreSQL (with replication)
```

**Considerations**:
- Session stickiness for Chainlit
- Shared session store (Redis) for agents
- Database connection pooling
- Rate limiting per user

## Monitoring & Observability

### Health Checks

Each service exposes health endpoints:
- Chainlit: `http://localhost:8082/`
- Agents: `http://localhost:8083/api/v1/health`
- MCP: `http://localhost:8081/`

### Logging

View logs:
```bash
docker compose logs -f chainlit agents mcp-server
```

### Metrics (Future)

- Request count and latency
- Tool call frequency
- Error rates
- Token usage (OpenAI)
- Database query performance

## Cost Estimation

### OpenAI API Costs (GPT-4)

- Input: $0.03 per 1K tokens
- Output: $0.06 per 1K tokens
- Average query: ~2K tokens total
- Cost per query: ~$0.10

**Example monthly usage**:
- 100 queries/day × 30 days = 3,000 queries
- Cost: ~$300/month

**Cost reduction**:
- Use GPT-3.5-turbo: ~$0.002 per query → ~$6/month
- Cache common queries
- Optimize prompts to reduce tokens

### Infrastructure Costs

**Docker on local machine**: Free

**Cloud hosting (AWS example)**:
- t3.medium (2 vCPU, 4GB RAM): ~$30/month
- RDS PostgreSQL (db.t3.micro): ~$15/month
- Total: ~$45/month (excluding API costs)

## Development Workflow

### Making Changes

1. **Update code** in `agents/` or `chainlit/`
2. **Rebuild container**:
   ```bash
   docker compose build agents
   ```
3. **Restart service**:
   ```bash
   docker compose up -d agents
   ```
4. **Test changes**:
   ```bash
   ./test-chat-services.sh
   ```

### Adding New Tools

1. **Add tool to MCP server** (Go code)
2. **Agent automatically discovers** new tool via MCP
3. **Update system prompt** if needed for context
4. **Test new functionality**

### Debugging

1. **Check logs**:
   ```bash
   docker compose logs -f agents
   ```

2. **Test API directly**:
   ```bash
   curl -X POST http://localhost:8083/api/v1/chat \
     -H "Content-Type: application/json" \
     -d '{"message": "test", "session_id": "debug"}'
   ```

3. **Connect to container**:
   ```bash
   docker compose exec agents bash
   ```

## Future Enhancements

### Short Term
- [ ] Streaming responses for better UX
- [ ] Chat history persistence
- [ ] Multiple agent types (budget, tax, investment)
- [ ] Error recovery and retry logic

### Medium Term
- [ ] User authentication and multi-user support
- [ ] File upload (CSV, bank statements)
- [ ] Export conversations
- [ ] Voice input/output
- [ ] Mobile app

### Long Term
- [ ] Multi-language support
- [ ] Custom agent training
- [ ] Integration with banking APIs
- [ ] Automated financial advice
- [ ] Compliance and audit logs
