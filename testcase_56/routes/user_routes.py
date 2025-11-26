from fastapi import APIRouter, Depends, HTTPException, status, Request
from models.user import NewEmailRequest, AuthUser, User
from services.user_service import user_service
from typing import Annotated

router = APIRouter(prefix="/profile", tags=["user"])

# --- Dependencies ---

def get_current_user(request: Request) -> AuthUser:
    """
    Simulates authentication by checking a valid session.
    
    SECURITY NOTE: The actual session validation and CSRF token enforcement 
    are handled by the middleware in main.py. This function assumes a valid 
    session token (session_id cookie) is present and maps to a user.
    
    In a real system, this would decode a JWT or check a session store based on the cookie.
    """
    # Placeholder: Assuming user ID 1 is authenticated for any valid session
    if not request.cookies.get("session_id"):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Not authenticated")
        
    return AuthUser(user_id=1)

# --- Routes ---

@router.get("/", response_model=User)
async def get_user_profile(current_user: Annotated[AuthUser, Depends(get_current_user)]):
    """Retrieves the authenticated user's profile."""
    user = user_service.get_user_by_id(current_user.user_id)
    if not user:
        # Should only happen if DB state is inconsistent
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
    return user

@router.post("/change_email", response_model=User, status_code=status.HTTP_200_OK)
async def change_user_email(
    request: Request,
    email_data: NewEmailRequest, 
    current_user: Annotated[AuthUser, Depends(get_current_user)]
):
    """
    Updates the user's email address.
    
    SECURITY DEFENSE: This state-changing endpoint is protected against CSRF 
    by the global CSRFMiddleware, which requires a valid X-CSRF-Token header.
    
    Input validation (email format) is handled automatically by Pydantic (NewEmailRequest).
    """
    
    new_email = email_data.new_email
    
    updated_user = user_service.update_user_email(
        user_id=current_user.user_id, 
        new_email=new_email
    )
    
    if not updated_user:
        # Implement proper error handling without leaking internal details
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Email update failed. User may not exist or email is invalid.")
        
    return updated_user