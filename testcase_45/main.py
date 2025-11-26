from fastapi import FastAPI
from api.v1 import calculator
import logging

# Configure basic logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Application metadata
TITLE = "Secure Expression Calculator API"
VERSION = "1.0.0"
DESCRIPTION = "A RESTful service for safely evaluating mathematical expressions, strictly preventing Remote Code Execution (RCE) via restricted AST parsing."

app = FastAPI(
    title=TITLE,
    version=VERSION,
    description=DESCRIPTION,
    docs_url="/docs",
    redoc_url="/redoc"
)

# --- API Routers ---
app.include_router(calculator.router, prefix="/api/v1")

@app.get("/", include_in_schema=False)
async def root():
    """Root endpoint for basic health check."""
    return {"message": "Welcome to the Secure Calculator API. Access /docs for API documentation."}

if __name__ == "__main__":
    # Standard security practice: In production, use a dedicated ASGI server (like gunicorn + uvicorn workers).
    # This block is for development/testing only.
    import uvicorn
    logger.info("Starting application in development mode.")
    uvicorn.run(app, host="0.0.0.0", port=8000)