import os
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from jose import JWTError, jwt
from typing import Annotated

from models.config import User

# SECURITY WARNING: In production, these must be loaded from secure environment variables or a secrets manager.
SECRET_KEY = "09d25e094f78a2e1c3b4f6d8c7a9b0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7"
ALGORITHM = "HS256"

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

def decode_access_token(token: str):
    """Decodes the JWT token and returns the payload, raising exceptions on failure."""
    try:
        # Ensure proper algorithm is used
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        
        # Input validation on payload structure
        username: str = payload.get("sub")
        user_id: int = payload.get("user_id")
        role: str = payload.get("role")

        if not all([username, user_id, role]):
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid token payload (missing required user data)",
                headers={"WWW-Authenticate": "Bearer"},
            )
        return {"username": username, "user_id": user_id, "role": role}
    except JWTError:
        # Generic error for token validation failure
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )

def get_current_user(token: Annotated[str, Depends(oauth2_scheme)]) -> User:
    """Dependency to get the current authenticated user from the JWT token."""
    try:
        payload = decode_access_token(token)
        # In a real application, you would verify the user still exists in the DB here
        user_data = User(
            username=payload["username"],
            user_id=payload["user_id"],
            role=payload["role"]
        )
        return user_data
    except HTTPException as e:
        raise e

def get_current_admin_user(current_user: Annotated[User, Depends(get_current_user)]) -> User:
    """Dependency to ensure the authenticated user has administrative privileges."""
    # Authorization Check (Least Privilege Principle)
    if current_user.role != "admin":
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Operation requires administrator privileges."
        )
    return current_user
