# Quick Start: GnuCash Chat Interface

Get started with the chat interface in 5 minutes.

## Prerequisites

- Docker and Docker Compose installed
- OpenAI API key (get from https://platform.openai.com/api-keys)
- GnuCash data in PostgreSQL database

## Step 1: Configure API Key

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` and add your OpenAI API key:
```bash
# Find this line and replace with your key
OPENAI_API_KEY=sk-your-actual-key-here
```

## Step 2: Start Services

Start all services with Docker Compose:
```bash
docker compose up -d
```

This will start:
- PostgreSQL database (port 5432)
- Backend API (port 8080)
- MCP Server (port 8081)
- Agents Service (port 8083)
- Chainlit UI (port 8082)

## Step 3: Verify Services

Run the test script:
```bash
./test-chat-services.sh
```

Or manually check:
```bash
# Check all services are running
docker compose ps

# Should show all services as "Up" and "healthy"
```

## Step 4: Open Chat Interface

Open your browser to:
```
http://localhost:8082
```

You should see the GnuCash AI Assistant welcome message.

## Step 5: Try It Out

Ask some questions:

```
What accounts do I have?
```

```
What's my checking account balance?
```

```
Show me expenses from last month
```

```
How much did I spend on groceries this year?
```

## Troubleshooting

### Services won't start

Check logs:
```bash
docker compose logs
```

### Chat interface shows error

1. Verify OpenAI API key is set:
```bash
grep OPENAI_API_KEY .env
```

2. Check agents service:
```bash
docker compose logs agents
```

3. Restart services:
```bash
docker compose restart agents chainlit
```

### Slow responses

- GPT-4 can take 10-30 seconds for complex queries
- This is normal behavior
- Consider using a faster model in `.env`:
```bash
MODEL_NAME=gpt-3.5-turbo
```

## Next Steps

- Read [README_CHAT.md](README_CHAT.md) for detailed documentation
- Explore [example queries](#example-queries)
- Check service logs: `docker compose logs -f chainlit agents`

## Example Queries

### Account Information
- "List all my accounts"
- "What's the hierarchy of my accounts?"
- "Show me all asset accounts"

### Balances
- "What's my net worth?"
- "Show balances for all accounts"
- "What's in my savings account?"

### Transactions
- "Show transactions from last week"
- "Find all transactions over $500"
- "What did I spend at Whole Foods?"

### Analysis
- "What are my top expense categories?"
- "How much did I earn this month?"
- "Show me spending trends"
- "What's my biggest expense?"

### Income & Expenses
- "What's my total income this year?"
- "Show me all expenses from December"
- "How much did I spend on utilities?"

## Stopping Services

Stop all services:
```bash
docker compose down
```

Stop specific services:
```bash
docker compose stop chainlit agents
```

## Getting Help

- View logs: `docker compose logs [service-name]`
- Check health: `./test-chat-services.sh`
- Read docs: [README_CHAT.md](README_CHAT.md)
- Report issues: GitHub Issues

## Configuration Options

### Change LLM Model

Edit `.env`:
```bash
# Faster, cheaper
MODEL_NAME=gpt-3.5-turbo

# More capable (default)
MODEL_NAME=gpt-4

# Latest model
MODEL_NAME=gpt-4-turbo
```

### Use Anthropic Claude

Edit `.env`:
```bash
MODEL_NAME=claude-3-sonnet-20240229
ANTHROPIC_API_KEY=your-anthropic-key
```

(Requires Strands to support Anthropic - check compatibility)

## Development Mode

Run services locally for development:

### Agents Service
```bash
cd agents
pip install -r requirements.txt
export OPENAI_API_KEY=sk-xxx
export MCP_SERVER_URL=http://localhost:8081
python main.py
```

### Chainlit UI
```bash
cd chainlit
pip install -r requirements.txt
export AGENTS_SERVICE_URL=http://localhost:8083
chainlit run app.py
```

## Architecture Overview

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ WebSocket
       ▼
┌─────────────┐
│  Chainlit   │ Port 8082 - Chat UI
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────┐
│   Agents    │ Port 8083 - AI Agents
└──────┬──────┘
       │ MCP Protocol
       ▼
┌─────────────┐
│ MCP Server  │ Port 8081 - Financial Data
└──────┬──────┘
       │ SQL
       ▼
┌─────────────┐
│ PostgreSQL  │ Port 5432 - Database
└─────────────┘
```

## What's Happening Under the Hood

1. **You type a message** in the Chainlit web interface
2. **Chainlit sends** your message to the Agents service via HTTP
3. **Agents service** uses a Strands Agent (powered by GPT-4) to understand your question
4. **The agent decides** which MCP tools to call (e.g., get_account_balance)
5. **MCP server** queries the PostgreSQL database
6. **Data flows back** through the chain to create a natural language response
7. **You see the answer** in the chat interface

## Common Issues

### "Error: OPENAI_API_KEY is required"
- Add your API key to `.env` file
- Restart services: `docker compose restart agents`

### "Agents service is not responding"
- Check if services are running: `docker compose ps`
- View logs: `docker compose logs agents`
- Restart: `docker compose restart agents`

### Port already in use
- Change ports in `docker-compose.yml`
- Or stop conflicting services

### Out of memory
- Increase Docker memory limit
- Check container resource usage: `docker stats`

## Support

- 📖 Documentation: [README_CHAT.md](README_CHAT.md)
- 🔧 Service docs: [agents/README.md](agents/README.md), [chainlit/README.md](chainlit/README.md)
- 🐛 Issues: GitHub Issues
- 💬 Community: GitHub Discussions
