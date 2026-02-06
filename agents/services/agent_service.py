from typing import Dict
from agents.finance_agent import FinanceAgent

class AgentService:
    """Service to manage agent instances per session."""

    def __init__(self):
        self.agents: Dict[str, FinanceAgent] = {}

    async def get_agent(self, session_id: str) -> FinanceAgent:
        """Get or create an agent for the session."""
        if session_id not in self.agents:
            agent = FinanceAgent()
            await agent.initialize()
            self.agents[session_id] = agent
        return self.agents[session_id]

    async def process_message(self, message: str, session_id: str) -> str:
        """Process a message using the session's agent."""
        agent = await self.get_agent(session_id)
        return await agent.query(message)

    def clear_session(self, session_id: str):
        """Clear an agent session."""
        if session_id in self.agents:
            # Cleanup MCP client context before removing
            self.agents[session_id].cleanup()
            del self.agents[session_id]
