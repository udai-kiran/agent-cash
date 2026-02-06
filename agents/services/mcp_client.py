from mcp.client.streamable_http import streamablehttp_client
from strands.tools.mcp import MCPClient
from config import settings

def create_mcp_client() -> MCPClient:
    """Create MCP client connected to GnuCash MCP server."""
    def create_transport():
        return streamablehttp_client(settings.MCP_SERVER_URL)

    return MCPClient(create_transport)
