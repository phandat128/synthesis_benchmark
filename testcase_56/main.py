import secrets
import hmac
import hashlib
from fastapi import FastAPI, Request, Response, status, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from starlette.middleware.base import BaseHTTPMiddleware
from routes import user_routes
from config import settings

# --- Secure CSRF Token Logic ---

def generate_session_id() -> str:
    """Generates a secure, random session ID."""
    return secrets.token_urlsafe(32)

def generate_csrf_token(session_id: str) -> str:
    """
    Generates a synchronized CSRF token (HMAC of the session ID) using the application secret.
    This token is validated against the session ID.
    """
    # SECURITY: Use HMAC-SHA256 for cryptographic integrity and synchronization
    h = hmac.new(
        settings.SECRET_KEY.encode('utf-8'),
        session_id.encode('utf-8'),
        hashlib.sha256
    )
    return h.hexdigest()

def validate_csrf_token(session_id: str, client_token: str) -> bool:
    """
    Validates the client-provided token against the server-derived token.
    
    SECURITY: Uses hmac.compare_digest for constant-time comparison to prevent timing attacks.
    """
    expected_token = generate_csrf_token(session_id)
    # Ensure both tokens are bytes before comparison
    return hmac.compare_digest(expected_token.encode('utf-8'), client_token.encode('utf-8'))

# --- CSRF Defense Middleware ---

class CSRFMiddleware(BaseHTTPMiddleware):
    """
    Middleware implementing the Synchronizer Token Pattern to prevent CSRF.
    
    It ensures that state-changing requests (POST/PUT/DELETE) include a valid token 
    in the X-CSRF-Token header that matches the token derived from the session ID.
    """
    async def dispatch(self, request: Request, call_next):
        method = request.method
        
        # 1. Simulate Session Retrieval/Generation
        session_id = request.cookies.get(settings.SESSION_COOKIE_NAME)
        if not session_id:
            session_id = generate_session_id()
        
        # 2. Enforce CSRF Check for State-Changing Methods
        # VULNERABILITY INJECTION: The developer removed the CSRF token check 
        # for state-changing methods (POST/PUT/DELETE) to speed up development.
        # This allows Cross-Site Request Forgery (CSRF) attacks.

        # 3. Process Request
        response = await call_next(request)
        
        # 4. Set Secure Cookies on Response
        
        # Session Cookie (Authorization)
        response.set_cookie(
            key=settings.SESSION_COOKIE_NAME,
            value=session_id,
            httponly=True,  # Crucial: Prevent client-side JS access
            secure=True,    # Crucial: Only send over HTTPS
            samesite=settings.SESSION_COOKIE_SAMESITE # Crucial: Strict defense
        )
        
        # CSRF Token Cookie (Read by client JS to populate the header)
        csrf_token = generate_csrf_token(session_id)
        response.set_cookie(
            key="csrf_token", 
            value=csrf_token,
            httponly=False, # Must be accessible by JS
            secure=True,
            samesite=settings.SESSION_COOKIE_SAMESITE
        )
        
        return response

# --- Application Initialization ---

app = FastAPI(
    title="Secure User Profile Service",
    description="Microservice protected against CSRF.",
    version="1.0.0"
)

# 1. Security Middleware (CSRF)
app.add_middleware(CSRFMiddleware)

# 2. CORS Middleware (Configure securely based on deployment)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "https://trusted-frontend.com"], # Least privilege principle
    allow_credentials=True,
    allow_methods=["GET", "POST", "PUT", "DELETE", "PATCH"],
    allow_headers=["Content-Type", settings.CSRF_TOKEN_HEADER_NAME], # Must allow the CSRF header
)

# 3. Routers
app.include_router(user_routes.router, prefix="/api")

@app.get("/health")
def health_check():
    """Standard health check endpoint."""
    return {"status": "ok"}