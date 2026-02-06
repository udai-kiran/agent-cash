# GnuCash Chat Interface - Quick Reference

## Essential Commands

### Start/Stop

```bash
# Start all services
docker compose up -d

# Start specific services
docker compose up -d chainlit agents

# Stop all services
docker compose down

# Stop specific service
docker compose stop chainlit

# Restart service
docker compose restart agents
```

### Testing

```bash
# Run automated tests
./test-chat-services.sh

# Test agents API
curl -X POST http://localhost:8083/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello", "session_id": "test"}'

# Check health
curl http://localhost:8082  # Chainlit
curl http://localhost:8083/api/v1/health  # Agents
curl http://localhost:8081  # MCP Server
```

### Monitoring

```bash
# View logs (follow mode)
docker compose logs -f chainlit agents

# View recent logs
docker compose logs --tail=50 agents

# Check service status
docker compose ps

# Check container resources
docker stats
```

### Troubleshooting

```bash
# Rebuild and restart
docker compose build agents chainlit
docker compose up -d agents chainlit

# Check environment variables
docker compose exec agents env | grep OPENAI_API_KEY

# Connect to container
docker compose exec agents bash

# View full logs
docker compose logs agents > agents.log
```

## Ports

| Service   | Port | Purpose                |
|-----------|------|------------------------|
| Chainlit  | 8082 | Chat UI (user-facing)  |
| Agents    | 8083 | AI agents API          |
| MCP       | 8081 | Financial data tools   |
| Backend   | 8080 | REST API (existing)    |
| PostgreSQL| 5432 | Database               |

## URLs

- Chat Interface: http://localhost:8082
- Agents API: http://localhost:8083/docs
- Backend API: http://localhost:8080/api/v1
- MCP Server: http://localhost:8081

## Environment Variables

```bash
# Required
OPENAI_API_KEY=sk-your-key-here

# Optional
MODEL_NAME=gpt-4  # or gpt-3.5-turbo
MCP_SERVER_URL=http://mcp-server:8081
AGENTS_SERVICE_URL=http://agents:8083
```

## Example Queries

### Accounts
```
List all my accounts
What accounts do I have?
Show me the account hierarchy
Search for accounts with "savings" in the name
```

### Balances
```
What's my checking account balance?
Show all account balances
What's my net worth?
What's the balance of my savings account?
```

### Transactions
```
Show transactions from last month
Find all transactions over $500
What did I spend at Whole Foods?
Show me recent transactions
```

### Analysis
```
What are my top 5 expense categories?
How much did I spend on groceries this month?
What's my total income this year?
Show me spending trends
Analyze my expenses
```

### Income & Expenses
```
Calculate my income and expenses for this month
What's my total income this year?
How much did I spend on utilities?
Show me all expenses from December
```

## File Structure

```
gnucash/
├── agents/                     # AI Agents Service
│   ├── agents/
│   │   └── finance_agent.py
│   ├── services/
│   │   ├── mcp_client.py
│   │   └── agent_service.py
│   ├── routes/
│   │   └── chat.py
│   ├── main.py
│   ├── config.py
│   └── Dockerfile
│
├── chainlit/                   # Chat UI Service
│   ├── utils/
│   │   └── agents_client.py
│   ├── .chainlit/
│   │   └── config.toml
│   ├── app.py
│   └── Dockerfile
│
├── backend/                    # Existing Backend
├── frontend/                   # Existing Frontend
├── docs/                       # Documentation
│   └── chat-architecture.md
│
├── docker-compose.yml          # All services
├── .env.example                # Config template
├── README_CHAT.md              # Full documentation
├── QUICKSTART_CHAT.md          # Quick start guide
└── test-chat-services.sh       # Test script
```

## Common Issues

### Issue: "Agents service is not responding"

**Solution:**
```bash
# Check if OpenAI API key is set
grep OPENAI_API_KEY .env

# View logs
docker compose logs agents

# Restart service
docker compose restart agents
```

### Issue: "Slow responses"

**Cause:** GPT-4 takes 10-30 seconds

**Solutions:**
1. Use faster model: `MODEL_NAME=gpt-3.5-turbo` in .env
2. Check MCP server: `curl http://localhost:8081`
3. Review logs: `docker compose logs mcp-server`

### Issue: Port already in use

**Solution:**
```bash
# Find what's using the port
lsof -i :8082

# Kill the process or change port in docker-compose.yml
```

### Issue: Services won't start

**Solution:**
```bash
# Check all services
docker compose ps

# View logs
docker compose logs

# Restart all
docker compose down
docker compose up -d
```

## API Endpoints

### Agents Service (Port 8083)

```bash
# Chat
POST /api/v1/chat
Body: {"message": "...", "session_id": "..."}

# Health
GET /api/v1/health

# Clear session
DELETE /api/v1/session/{session_id}

# API Docs
GET /docs
```

## MCP Tools Available

1. `get_accounts` - List all accounts
2. `get_account_by_id` - Get specific account
3. `get_account_balance` - Get account balance
4. `get_transactions` - Query transactions
5. `get_transaction_by_id` - Get specific transaction
6. `calculate_income_expenses` - Calculate totals
7. `get_account_hierarchy` - Get account tree
8. `search_accounts` - Search accounts by name
9. `search_transactions` - Search transactions
10. `get_commodity_prices` - Get currency prices
11. `get_budget_info` - Get budget data

## Configuration Files

### .env
```bash
OPENAI_API_KEY=sk-your-key-here
MODEL_NAME=gpt-4
```

### agents/config.py
```python
class Settings(BaseSettings):
    HOST: str = "0.0.0.0"
    PORT: int = 8083
    MCP_SERVER_URL: str = "http://mcp-server:8081"
    OPENAI_API_KEY: str
    MODEL_NAME: str = "gpt-4"
```

### chainlit/.chainlit/config.toml
```toml
[project]
enable_telemetry = false
session_timeout = 3600

[UI]
name = "GnuCash AI Assistant"
```

## Development

### Local Development (without Docker)

**Agents service:**
```bash
cd agents
pip install -r requirements.txt
export OPENAI_API_KEY=sk-xxx
export MCP_SERVER_URL=http://localhost:8081
python main.py
```

**Chainlit:**
```bash
cd chainlit
pip install -r requirements.txt
export AGENTS_SERVICE_URL=http://localhost:8083
chainlit run app.py
```

### Adding New Features

**New agent type:**
1. Create `agents/agents/new_agent.py`
2. Update `services/agent_service.py`
3. Add route if needed

**New UI feature:**
1. Edit `chainlit/app.py`
2. Rebuild: `docker compose build chainlit`
3. Restart: `docker compose up -d chainlit`

## Performance Tips

1. **Use GPT-3.5-turbo** for faster responses
2. **Cache common queries** in Redis
3. **Parallel tool calls** where possible
4. **Optimize database queries** in MCP server
5. **Add response streaming** for better UX

## Security Best Practices

1. ✅ Store API keys in .env (never commit)
2. ✅ Use internal Docker network
3. ⚠️ Add user authentication (future)
4. ⚠️ Restrict CORS origins (production)
5. ⚠️ Implement rate limiting (production)

## Cost Management

### OpenAI API Costs

| Model | Input | Output | Avg Query |
|-------|-------|--------|-----------|
| GPT-4 | $0.03/1K | $0.06/1K | ~$0.10 |
| GPT-3.5-turbo | $0.0005/1K | $0.0015/1K | ~$0.002 |

### Monthly Estimates (100 queries/day)

- **GPT-4**: ~$300/month
- **GPT-3.5-turbo**: ~$6/month

### Reduce Costs

1. Use GPT-3.5-turbo
2. Implement caching
3. Optimize prompts
4. Set usage limits

## Resources

### Documentation
- 📘 Full Guide: [README_CHAT.md](README_CHAT.md)
- 🚀 Quick Start: [QUICKSTART_CHAT.md](QUICKSTART_CHAT.md)
- 🏗️ Architecture: [docs/chat-architecture.md](docs/chat-architecture.md)
- 📋 Summary: [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)

### Service Docs
- 🤖 Agents: [agents/README.md](agents/README.md)
- 💬 Chainlit: [chainlit/README.md](chainlit/README.md)

### External Resources
- Chainlit Docs: https://docs.chainlit.io
- Strands Agents: https://docs.strands.ai
- MCP Spec: https://modelcontextprotocol.io
- OpenAI API: https://platform.openai.com/docs

## Support

- 🐛 Report Issues: GitHub Issues
- 💬 Discussions: GitHub Discussions
- 📧 Email: support@example.com
- 📝 Logs: `docker compose logs [service]`

## Version Info

- **Chainlit**: 2.9.6+
- **Strands**: 0.4.0
- **FastAPI**: 0.115.0
- **MCP**: 1.2.0
- **Python**: 3.11

---

**Last Updated:** 2026-02-05
**Version:** 1.0.0
