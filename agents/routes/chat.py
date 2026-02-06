from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from services.agent_service import AgentService

router = APIRouter()
agent_service = AgentService()

class ChatRequest(BaseModel):
    message: str
    session_id: str

class ChatResponse(BaseModel):
    response: str
    session_id: str

@router.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    """Handle chat message from Chainlit."""
    try:
        response = await agent_service.process_message(
            request.message,
            request.session_id
        )
        return ChatResponse(
            response=response,
            session_id=request.session_id
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/health")
async def health():
    """Health check endpoint."""
    return {"status": "healthy"}

@router.delete("/session/{session_id}")
async def clear_session(session_id: str):
    """Clear a session's agent."""
    agent_service.clear_session(session_id)
    return {"status": "cleared"}
