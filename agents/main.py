from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from routes import chat
from config import settings

app = FastAPI(
    title="GnuCash Agents Service",
    description="AI agents for financial data analysis",
    version="1.0.0"
)

# CORS configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routes
app.include_router(chat.router, prefix="/api/v1", tags=["chat"])

@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "service": "GnuCash Agents",
        "version": "1.0.0",
        "status": "running"
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=False
    )
