from pydantic import BaseModel, EmailStr, Field

class User(BaseModel):
    """Schema for internal user representation (simulated DB model)"""
    id: int
    username: str
    email: EmailStr
    
class NewEmailRequest(BaseModel):
    """Input schema for changing the user's email"""
    # SECURITY: Pydantic ensures the input is a valid email format, preventing basic injection/malformed data.
    new_email: EmailStr = Field(..., description="The new email address for the user.")

class AuthUser(BaseModel):
    """Schema representing the authenticated user ID (used in dependencies)"""
    user_id: int