# GnuCash Agents Service

AI agent service that bridges Chainlit chat interface with the MCP server.

## Overview

This service provides a REST API that processes natural language queries about financial data. It uses Strands Agents framework to orchestrate AI models with access to GnuCash data through the MCP server.

## Architecture

```
FastAPI Server
    ↓
Agent Service (session management)
    ↓
Finance Agent (Strands)
    ↓
MCP Client
    ↓
MCP Server (port 8081)
```

## Components

### Routes (`routes/chat.py`)
- `POST /api/v1/chat` - Process chat messages
- `GET /api/v1/health` - Health check
- `DELETE /api/v1/session/{session_id}` - Clear session

### Services
- `agent_service.py` - Manages agent instances per session
- `mcp_client.py` - MCP server connection

### Agents
- `finance_agent.py` - Main financial assistant agent

## Configuration

Environment variables (set in `.env` or docker-compose.yml):

```bash
HOST=0.0.0.0                          # Server host
PORT=8083                             # Server port
MCP_SERVER_URL=http://mcp-server:8081 # MCP server endpoint
OPENAI_API_KEY=sk-xxx                 # OpenAI API key
MODEL_NAME=gpt-4                      # LLM model to use
```

## API Usage

### Chat Endpoint

```bash
curl -X POST http://localhost:8083/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What is my checking account balance?",
    "session_id": "user-123"
  }'
```

Response:
```json
{
  "response": "Your checking account balance is $5,234.56",
  "session_id": "user-123"
}
```

### Health Check

```bash
curl http://localhost:8083/api/v1/health
```

## Local Development

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set environment variables:
```bash
export OPENAI_API_KEY=sk-your-key-here
export MCP_SERVER_URL=http://localhost:8081
```

3. Run server:
```bash
python main.py
```

4. Test endpoint:
```bash
curl -X POST http://localhost:8083/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello", "session_id": "test"}'
```

## Docker Development

Build and run:
```bash
docker compose up -d agents
```

View logs:
```bash
docker compose logs -f agents
```

Restart:
```bash
docker compose restart agents
```

## Adding New Agents

To add a specialized agent:

1. Create agent file in `agents/`:
```python
from strands import Agent
from services.mcp_client import create_mcp_client

class BudgetAgent:
    def __init__(self):
        self.mcp_client = create_mcp_client()
        self.agent = None

    async def initialize(self):
        async with self.mcp_client as client:
            tools = await client.list_tools()
            self.agent = Agent(
                tools=tools,
                model="gpt-4",
                system_prompt="You are a budget advisor..."
            )

    async def query(self, message: str) -> str:
        response = await self.agent.ainvoke(message)
        return response
```

2. Update `agent_service.py` to route to appropriate agent

3. Add route in `routes/chat.py` if needed

## Troubleshooting

### Agent not responding

Check MCP server connection:
```bash
curl http://localhost:8081
```

View detailed logs:
```bash
docker compose logs agents
```

### OpenAI API errors

Verify API key is set:
```bash
docker compose exec agents env | grep OPENAI_API_KEY
```

Check API key validity at https://platform.openai.com

### Slow responses

- GPT-4 typically takes 10-30 seconds
- Consider using gpt-3.5-turbo for faster responses
- Check MCP server response times

## Dependencies

- **FastAPI** - Web framework
- **Strands Agents** - AI agent orchestration
- **MCP SDK** - Model Context Protocol client
- **aiohttp** - Async HTTP client
- **pydantic** - Data validation

## Security

- API keys stored in environment variables
- No authentication currently implemented
- CORS allows all origins (configure for production)
- Internal Docker network for service communication

## Future Enhancements

- [ ] Add streaming response support
- [ ] Implement rate limiting
- [ ] Add user authentication
- [ ] Support multiple specialized agents
- [ ] Cache frequent queries
- [ ] Add metrics and monitoring
