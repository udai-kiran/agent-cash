import chainlit as cl
from utils.agents_client import AgentsClient
import uuid

agents_client = AgentsClient()

@cl.on_chat_start
async def start():
    """Initialize chat session."""
    # Generate unique session ID
    session_id = str(uuid.uuid4())
    cl.user_session.set("session_id", session_id)

    # Check if agents service is healthy
    is_healthy = await agents_client.health_check()
    if not is_healthy:
        await cl.Message(
            content="⚠️ Warning: Agents service is not responding. Please check if all services are running."
        ).send()
        return

    # Send welcome message
    await cl.Message(
        content="""👋 Welcome to GnuCash AI Assistant!

I can help you with:
- 📊 Checking account balances
- 💰 Analyzing income and expenses
- 📈 Viewing financial trends
- 🔍 Finding specific transactions
- 💡 Getting financial insights

What would you like to know about your finances?"""
    ).send()

@cl.on_message
async def main(message: cl.Message):
    """Handle incoming user messages."""
    session_id = cl.user_session.get("session_id")

    # Create a placeholder message
    msg = cl.Message(content="")
    await msg.send()

    try:
        # Show thinking indicator
        msg.content = "🤔 Thinking..."
        await msg.update()

        # Send to agents service
        response = await agents_client.send_message(
            message.content,
            session_id
        )

        # Update message with response
        msg.content = response
        await msg.update()

    except Exception as e:
        error_message = str(e)
        msg.content = f"""❌ Error communicating with the AI agent.

**Error details:** {error_message}

**Troubleshooting:**
1. Check if all services are running: `docker compose ps`
2. Check agents service logs: `docker compose logs agents`
3. Verify your OPENAI_API_KEY is set correctly

Please try again or contact support if the issue persists."""
        await msg.update()

@cl.on_chat_end
async def end():
    """Clean up when chat ends."""
    session_id = cl.user_session.get("session_id")
    if session_id:
        # Optionally notify agents service to clean up session
        # This is a fire-and-forget, don't wait for response
        pass
