import aiohttp
from typing import Optional
import os

class AgentsClient:
    """Client for communicating with agents service."""

    def __init__(self, base_url: Optional[str] = None):
        self.base_url = base_url or os.getenv("AGENTS_SERVICE_URL", "http://agents:8083")

    async def send_message(
        self,
        message: str,
        session_id: str
    ) -> str:
        """Send message to agents service and get response."""
        async with aiohttp.ClientSession() as session:
            async with session.post(
                f"{self.base_url}/api/v1/chat",
                json={
                    "message": message,
                    "session_id": session_id
                },
                timeout=aiohttp.ClientTimeout(total=120)  # 2 minute timeout for agent responses
            ) as response:
                if response.status == 200:
                    data = await response.json()
                    return data["response"]
                else:
                    error_text = await response.text()
                    raise Exception(f"Agent service error ({response.status}): {error_text}")

    async def health_check(self) -> bool:
        """Check if agents service is healthy."""
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(
                    f"{self.base_url}/api/v1/health",
                    timeout=aiohttp.ClientTimeout(total=5)
                ) as response:
                    return response.status == 200
        except:
            return False
