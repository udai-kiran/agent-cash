# GnuCash Chat Interface

## Quick Start

1. Set your OpenAI API key:
   ```bash
   export OPENAI_API_KEY=sk-your-key-here
   ```

   Or create a `.env` file (copy from `.env.example`):
   ```bash
   cp .env.example .env
   # Edit .env and add your API key
   ```

2. Start all services:
   ```bash
   docker compose up -d
   ```

3. Open chat interface:
   ```
   http://localhost:8082
   ```

## Architecture

The chat interface consists of three main services:

```
User Browser
    ↓ HTTP/WebSocket
Chainlit UI (port 8082) - Chat interface
    ↓ HTTP REST API
Agents Service (port 8083) - AI agents with financial tools
    ↓ MCP Protocol
MCP Server (port 8081) - Financial data access
    ↓ SQL
PostgreSQL Database
```

### Services

- **Chainlit UI** (port 8082) - User-facing chat interface
  - Built with Chainlit framework
  - Handles WebSocket connections
  - Manages user sessions

- **Agents Service** (port 8083) - AI agent orchestration
  - Built with FastAPI
  - Strands-based finance agents
  - Connects to MCP server for data access
  - REST API for chat interactions

- **MCP Server** (port 8081) - Financial data provider
  - 11 financial data tools
  - Direct PostgreSQL access
  - Model Context Protocol interface

## Example Queries

Try these queries in the chat interface:

- "What's my checking account balance?"
- "Show me expenses from last month"
- "How much did I spend on groceries?"
- "What's my total net worth?"
- "Find all transactions over $500"
- "What are my top 5 expense categories?"
- "Show me income for the current year"
- "List all my accounts"

## Configuration

### Environment Variables

Edit `.env` file to configure:

- `OPENAI_API_KEY` - Required for agents (get from https://platform.openai.com)
- `MODEL_NAME` - LLM model to use (default: gpt-4)

### Supported Models

The agents service uses Strands, which supports multiple LLM providers:

- **OpenAI**: gpt-4, gpt-4-turbo, gpt-3.5-turbo
- **Anthropic**: claude-3-opus, claude-3-sonnet
- **Others**: Check Strands documentation for full list

## Troubleshooting

### Check Service Status

```bash
# View all services
docker compose ps

# Should show all services as "healthy"
```

### View Logs

```bash
# Chainlit logs
docker compose logs chainlit

# Agents service logs
docker compose logs agents

# MCP server logs
docker compose logs mcp-server

# Follow logs in real-time
docker compose logs -f chainlit agents
```

### Health Checks

```bash
# Chainlit
curl http://localhost:8082

# Agents service
curl http://localhost:8083/api/v1/health

# MCP server
curl http://localhost:8081
```

### Common Issues

**Issue**: "Agents service is not responding"
- Check if OPENAI_API_KEY is set correctly
- Verify MCP server is running: `docker compose ps mcp-server`
- Check agents logs: `docker compose logs agents`

**Issue**: "Agent service error (500)"
- Check if OpenAI API key is valid
- Verify MCP server is accessible
- Check agents service logs for detailed error

**Issue**: Chat interface won't load
- Verify all services are running: `docker compose ps`
- Check if port 8082 is available
- View chainlit logs: `docker compose logs chainlit`

**Issue**: Slow responses
- GPT-4 can take 10-30 seconds for complex queries
- Check MCP server response times
- Consider using faster model (gpt-3.5-turbo)

### Restart Services

```bash
# Restart specific service
docker compose restart chainlit
docker compose restart agents

# Restart all services
docker compose restart

# Complete rebuild
docker compose down
docker compose up -d --build
```

## Development

### Project Structure

```
.
├── agents/                    # Agents service
│   ├── agents/               # Agent implementations
│   │   └── finance_agent.py
│   ├── services/             # Core services
│   │   ├── mcp_client.py    # MCP connection
│   │   └── agent_service.py # Agent orchestration
│   ├── routes/               # API routes
│   │   └── chat.py
│   ├── main.py               # FastAPI app
│   ├── config.py             # Configuration
│   └── Dockerfile
│
├── chainlit/                  # Chainlit UI service
│   ├── utils/
│   │   └── agents_client.py  # Agents service client
│   ├── .chainlit/
│   │   └── config.toml       # Chainlit config
│   ├── app.py                # Chainlit app
│   └── Dockerfile
│
└── docker-compose.yml         # Service orchestration
```

### Local Development

Run agents service locally:
```bash
cd agents
pip install -r requirements.txt
export OPENAI_API_KEY=sk-your-key-here
export MCP_SERVER_URL=http://localhost:8081
python main.py
```

Run Chainlit locally:
```bash
cd chainlit
pip install -r requirements.txt
export AGENTS_SERVICE_URL=http://localhost:8083
chainlit run app.py
```

### API Documentation

Agents service provides OpenAPI documentation:
```
http://localhost:8083/docs
```

### Testing the Agent Service

Test chat endpoint directly:
```bash
curl -X POST http://localhost:8083/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What accounts exist?",
    "session_id": "test-123"
  }'
```

## Security Considerations

1. **API Key Management**:
   - Never commit `.env` file with real API keys
   - Use environment variables in production
   - Rotate keys regularly

2. **Network Isolation**:
   - Services communicate via internal Docker network
   - Only necessary ports exposed to host

3. **CORS**:
   - Currently allows all origins (*)
   - Configure for production use

4. **Authentication**:
   - Consider adding user authentication
   - Implement rate limiting for production

## Future Enhancements

Planned features:

1. **User Authentication**
   - Login/signup system
   - Per-user session management
   - Multi-user support

2. **Chat History**
   - Persist conversations in PostgreSQL
   - Chat history browser
   - Export conversations

3. **Advanced Agents**
   - Budget advisor agent
   - Tax assistant agent
   - Investment analyzer agent

4. **Streaming Responses**
   - Real-time response streaming
   - Better user experience for long queries

5. **File Upload**
   - Import transactions from CSV
   - Parse bank statements
   - Receipt OCR

6. **Dashboards**
   - Visual analytics
   - Interactive charts
   - Financial insights

## Contributing

To add new features:

1. Add tools to MCP server (in `backend/internal/infrastructure/mcp/`)
2. Update finance agent prompt if needed
3. Test with Chainlit interface
4. Submit pull request

## Support

- Report issues: https://github.com/your-repo/gnucash/issues
- Documentation: See README.md files in each service directory
- MCP Server docs: See backend/README_MCP.md

## License

Same as main GnuCash project
