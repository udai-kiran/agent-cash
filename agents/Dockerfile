FROM python:3.11-slim

WORKDIR /app

# Copy only pyproject.toml first for caching
COPY pyproject.toml .

# Install dependencies using pip
RUN pip install --no-cache-dir \
    "aiohttp>=3.13.3" \
    "fastapi>=0.128.1" \
    "mcp>=1.26.0" \
    "openai>=1.0.0" \
    "pydantic-settings>=2.12.0" \
    "strands-agents>=1.24.0" \
    "strands-agents-tools>=0.2.19" \
    "uvicorn[standard]>=0.30.0"

# Copy application code (but not .venv due to .dockerignore)
COPY . .

# Create non-root user
RUN useradd -m -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8083

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD python -c "import requests; requests.get('http://localhost:8083/api/v1/health')"

CMD ["python", "main.py"]
