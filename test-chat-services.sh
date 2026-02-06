#!/bin/bash

# Test script for GnuCash Chat Services
# This script verifies all services are running and accessible

set -e

echo "========================================="
echo "GnuCash Chat Services Test Script"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print success
success() {
    echo -e "${GREEN}✓${NC} $1"
}

# Function to print error
error() {
    echo -e "${RED}✗${NC} $1"
}

# Function to print warning
warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if services are running
echo "1. Checking Docker services..."
echo ""

if ! docker compose ps | grep -q "gnucash-mcp-server"; then
    error "MCP server is not running"
    exit 1
fi
success "MCP server container exists"

if ! docker compose ps | grep -q "gnucash-agents"; then
    error "Agents service is not running"
    exit 1
fi
success "Agents service container exists"

if ! docker compose ps | grep -q "gnucash-chainlit"; then
    error "Chainlit service is not running"
    exit 1
fi
success "Chainlit service container exists"

echo ""
echo "2. Testing service health endpoints..."
echo ""

# Test MCP server
if curl -s -f http://localhost:8081 > /dev/null 2>&1; then
    success "MCP server is responding (port 8081)"
else
    error "MCP server is not responding"
    exit 1
fi

# Test Agents service
if curl -s -f http://localhost:8083/api/v1/health > /dev/null 2>&1; then
    success "Agents service is responding (port 8083)"
else
    error "Agents service is not responding"
    warning "Check logs: docker compose logs agents"
    exit 1
fi

# Test Chainlit
if curl -s -f http://localhost:8082 > /dev/null 2>&1; then
    success "Chainlit is responding (port 8082)"
else
    error "Chainlit is not responding"
    warning "Check logs: docker compose logs chainlit"
    exit 1
fi

echo ""
echo "3. Testing Agents service API..."
echo ""

# Test chat endpoint
RESPONSE=$(curl -s -X POST http://localhost:8083/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello", "session_id": "test-'$(date +%s)'"}' \
  2>&1)

if echo "$RESPONSE" | grep -q "response"; then
    success "Chat endpoint is working"
    echo "   Response preview: $(echo $RESPONSE | head -c 100)..."
else
    error "Chat endpoint returned unexpected response"
    echo "   Response: $RESPONSE"
    warning "This might be normal if OPENAI_API_KEY is not set"
fi

echo ""
echo "4. Checking environment configuration..."
echo ""

if [ -f .env ]; then
    success ".env file exists"

    if grep -q "OPENAI_API_KEY=sk-" .env 2>/dev/null; then
        success "OPENAI_API_KEY appears to be set"
    else
        warning "OPENAI_API_KEY may not be set correctly"
        echo "   Set it in .env file: OPENAI_API_KEY=sk-your-key-here"
    fi
else
    warning ".env file not found"
    echo "   Copy from .env.example and add your OpenAI API key"
fi

echo ""
echo "5. Checking service logs for errors..."
echo ""

# Check for errors in agents service
AGENT_ERRORS=$(docker compose logs agents --tail=50 2>&1 | grep -i "error" | wc -l)
if [ "$AGENT_ERRORS" -eq 0 ]; then
    success "No errors in agents service logs"
else
    warning "Found $AGENT_ERRORS error lines in agents service logs"
    echo "   View logs: docker compose logs agents"
fi

# Check for errors in chainlit service
CHAINLIT_ERRORS=$(docker compose logs chainlit --tail=50 2>&1 | grep -i "error" | wc -l)
if [ "$CHAINLIT_ERRORS" -eq 0 ]; then
    success "No errors in chainlit service logs"
else
    warning "Found $CHAINLIT_ERRORS error lines in chainlit service logs"
    echo "   View logs: docker compose logs chainlit"
fi

echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
echo ""
echo "All critical services are running!"
echo ""
echo "Next steps:"
echo "1. Open browser to http://localhost:8082"
echo "2. Try asking: 'What accounts exist?'"
echo "3. Try asking: 'Show me my account balances'"
echo ""
echo "Useful commands:"
echo "  docker compose logs -f chainlit agents  # Follow logs"
echo "  docker compose restart agents           # Restart agents"
echo "  docker compose ps                       # Service status"
echo ""
