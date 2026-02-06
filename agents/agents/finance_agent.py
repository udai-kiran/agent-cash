from strands import Agent
from strands.models import OpenAIModel
from services.mcp_client import create_mcp_client
from config import settings
from openai import AsyncOpenAI

class FinanceAgent:
    """Finance agent with access to GnuCash data."""

    def __init__(self):
        self.mcp_client = create_mcp_client()
        self.agent = None
        self.mcp_context = None
        self.openai_client = None

    async def initialize(self):
        """Initialize agent with MCP tools."""
        # Enter MCP client context and keep it open
        self.mcp_context = self.mcp_client.__enter__()

        tools = self.mcp_client.list_tools_sync()

        # Configure OpenAI model with custom API URL if needed
        if settings.OPENAI_API_URL != "https://api.openai.com/v1":
            # Custom API endpoint (e.g., OpenRouter)
            # Create OpenAI client with default headers for OpenRouter
            self.openai_client = AsyncOpenAI(
                api_key=settings.OPENAI_API_KEY,
                base_url=settings.OPENAI_API_URL,
                default_headers={
                    "HTTP-Referer": "https://gnucash-agents.local",
                    "X-Title": "GnuCash AI Assistant"
                }
            )

            model = OpenAIModel(
                client=self.openai_client,
                model_id=settings.MODEL_NAME,
                params={"temperature": 0.7}
            )
        else:
            # Standard OpenAI
            model = OpenAIModel(
                model_id=settings.MODEL_NAME,
                client_args={"api_key": settings.OPENAI_API_KEY},
                params={"temperature": 0.7}
            )

        # Create agent with tools and model
        self.agent = Agent(
            tools=tools,
            model=model,
            system_prompt="""You are a financial assistant with access to GnuCash data.
            You can help users understand their finances, analyze spending, track income,
            and answer questions about their accounts and transactions.

            When users ask questions:
            1. Use the available tools to fetch accurate data from GnuCash
            2. Present information clearly with relevant numbers and formatting
            3. Provide helpful insights and context
            4. Be conversational and friendly

            Available tools allow you to:
            - Get account balances and details
            - Query transactions with various filters
            - Calculate income and expenses
            - Analyze spending by category
            - Track account hierarchies
            """
        )

    def cleanup(self):
        """Cleanup MCP client context."""
        if self.mcp_context is not None:
            self.mcp_client.__exit__(None, None, None)
            self.mcp_context = None

    async def query(self, message: str) -> str:
        """Process user query and return response."""
        if not self.agent:
            await self.initialize()

        # Strands agents use synchronous invocation (agent is callable)
        result = self.agent(message)

        # Extract text from AgentResult
        # The result has a message attribute with content array
        if hasattr(result, 'message') and result.message:
            # Extract text from message content
            content_parts = []
            for content in result.message.get('content', []):
                if isinstance(content, dict) and 'text' in content:
                    content_parts.append(content['text'])
                elif isinstance(content, str):
                    content_parts.append(content)
            return '\n'.join(content_parts) if content_parts else str(result)

        # Fallback to string conversion
        return str(result)
