# GnuCash Chainlit UI

Chat interface for interacting with GnuCash financial data.

## Overview

This service provides a web-based chat interface using the Chainlit framework. Users can ask natural language questions about their finances, and the interface communicates with the Agents service to get responses.

## Architecture

```
User Browser (WebSocket)
    ↓
Chainlit App (FastAPI + React)
    ↓
Agents Client (HTTP)
    ↓
Agents Service (port 8083)
```

## Components

### Main Application (`app.py`)

- `@cl.on_chat_start` - Initialize new chat sessions
- `@cl.on_message` - Handle incoming messages
- `@cl.on_chat_end` - Clean up sessions

### Utils

- `agents_client.py` - HTTP client for agents service

### Configuration (`.chainlit/config.toml`)

Chainlit framework configuration:
- UI settings
- Session timeouts
- Feature flags

## Usage

### Starting the Service

With Docker:
```bash
docker compose up -d chainlit
```

Standalone:
```bash
chainlit run app.py --host 0.0.0.0 --port 8082
```

### Accessing the Interface

Open browser to:
```
http://localhost:8082
```

### Example Conversations

**User**: "What's my checking account balance?"

**Assistant**: "Your Assets:Current Assets:Checking Account has a balance of $5,234.56 as of today."

**User**: "Show me expenses from last month"

**Assistant**: "Here are your expenses from last month:
- Groceries: $456.78
- Utilities: $234.50
- Transportation: $189.23
Total: $880.51"

## Configuration

Environment variables:

```bash
AGENTS_SERVICE_URL=http://agents:8083  # Agents service endpoint
```

Chainlit config (`.chainlit/config.toml`):

```toml
[project]
enable_telemetry = false
session_timeout = 3600

[UI]
name = "GnuCash AI Assistant"
```

## Local Development

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set environment:
```bash
export AGENTS_SERVICE_URL=http://localhost:8083
```

3. Run Chainlit:
```bash
chainlit run app.py
```

4. Open browser to `http://localhost:8000`

## Customization

### Changing UI Theme

Edit `.chainlit/config.toml`:

```toml
[UI]
name = "My Financial Assistant"
description = "Custom description"
theme = "light"  # or "dark"
```

### Adding Custom Welcome Message

Edit `app.py`:

```python
@cl.on_chat_start
async def start():
    await cl.Message(
        content="Your custom welcome message here"
    ).send()
```

### Adding Message History

Store messages in PostgreSQL:

```python
@cl.on_message
async def main(message: cl.Message):
    # Store message
    await store_message(message.content)

    # Get response
    response = await agents_client.send_message(...)

    # Store response
    await store_message(response, is_assistant=True)
```

## Docker Development

Build:
```bash
docker compose build chainlit
```

Run:
```bash
docker compose up -d chainlit
```

Logs:
```bash
docker compose logs -f chainlit
```

Restart:
```bash
docker compose restart chainlit
```

## Troubleshooting

### Interface not loading

Check service status:
```bash
docker compose ps chainlit
```

View logs:
```bash
docker compose logs chainlit
```

Check port availability:
```bash
lsof -i :8082
```

### "Agents service is not responding"

1. Verify agents service is running:
```bash
curl http://localhost:8083/api/v1/health
```

2. Check network connectivity:
```bash
docker compose exec chainlit ping agents
```

3. Check environment variable:
```bash
docker compose exec chainlit env | grep AGENTS_SERVICE_URL
```

### Slow responses

- Agents service may take 10-30 seconds
- Check agents service logs for details
- Consider showing typing indicator

### WebSocket disconnections

- Check browser console for errors
- Verify firewall/proxy settings
- Increase session timeout in config.toml

## Chainlit Features

### File Upload

Enable file uploads (future enhancement):

```python
@cl.on_message
async def main(message: cl.Message):
    if message.elements:
        for element in message.elements:
            # Process uploaded file
            pass
```

### Actions

Add clickable actions:

```python
actions = [
    cl.Action(name="view_accounts", value="accounts", label="View Accounts"),
    cl.Action(name="view_expenses", value="expenses", label="View Expenses"),
]

await cl.Message(
    content="What would you like to do?",
    actions=actions
).send()
```

### Data Elements

Show structured data:

```python
await cl.Message(
    content="Here are your accounts:",
    elements=[
        cl.Text(name="accounts", content="Account data...")
    ]
).send()
```

## Dependencies

- **Chainlit** - Chat UI framework
- **aiohttp** - Async HTTP client
- **pydantic** - Data validation

## Security

- No authentication currently implemented
- All origins allowed (configure CORS)
- Consider adding:
  - User login system
  - Session encryption
  - Rate limiting

## Performance

- WebSocket for real-time communication
- Async operations throughout
- Connection pooling for agents client

## Future Enhancements

- [ ] User authentication
- [ ] Chat history persistence
- [ ] Export conversations
- [ ] File upload support
- [ ] Streaming responses
- [ ] Multi-language support
- [ ] Voice input/output
- [ ] Mobile app
