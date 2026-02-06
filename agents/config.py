from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # Server settings
    HOST: str = "0.0.0.0"
    PORT: int = 8083

    # MCP server connection
    MCP_SERVER_URL: str = "http://mcp-server:8081"

    # Agent configuration
    OPENAI_API_KEY: str
    OPENAI_API_URL: str = "https://api.openai.com/v1"
    MODEL_NAME: str = "gpt-4"

    class Config:
        env_file = ".env"

settings = Settings()
