FROM python:3.11-slim

WORKDIR /app

# Copy only pyproject.toml first for caching
COPY pyproject.toml .

# Install dependencies using pip
RUN pip install --no-cache-dir \
    "chainlit>=2.9.6" \
    "aiohttp>=3.10.0" \
    "pydantic>=2.10.0"

# Copy application (but not .venv due to .dockerignore)
COPY . .

# Create .chainlit directory and ensure proper permissions
RUN mkdir -p /app/.chainlit && \
    chown -R 1000:1000 /app

# Create non-root user
RUN useradd -m -u 1000 appuser

USER appuser

EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8082')"

# Set Chainlit environment variables to disable telemetry and configure
ENV CHAINLIT_TELEMETRY=false

CMD ["chainlit", "run", "app.py", "--host", "0.0.0.0", "--port", "8082"]
