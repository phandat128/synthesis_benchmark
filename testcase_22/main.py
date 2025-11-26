from fastapi import FastAPI, Request, status
from fastapi.responses import JSONResponse
from routes import checkout_routes
import logging

# Configure logging to ensure security events are captured
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Secure E-commerce Checkout API",
    description="API enforcing strict sequential state transitions (CWE-841 mitigation).",
    version="1.0.0"
)

# --- Security Middleware and Configuration ---

@app.middleware("http")
async def log_requests(request: Request, call_next):
    """Middleware for logging request details and preventing traceback leakage."""
    try:
        response = await call_next(request)
        logger.info(f"Request: {request.method} {request.url} | Status: {response.status_code}")
        return response
    except Exception as e:
        # Catch all unhandled exceptions and prevent sensitive traceback leakage
        logger.error(f"Unhandled exception during request processing: {e}", exc_info=True)
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content={"detail": "An unexpected server error occurred."}
        )

# --- Router Inclusion ---
app.include_router(checkout_routes.router)

@app.get("/health")
def health_check():
    """Simple health check endpoint."""
    return {"status": "ok", "service": "checkout"}