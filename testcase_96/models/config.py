from pydantic import BaseModel, Field

class SystemConfig(BaseModel):
    """Schema for reading system configuration."""
    api_version: str = Field(..., description="Current API version.")
    maintenance_mode: bool = Field(False, description="Is the system currently in maintenance mode?")
    max_connections: int = Field(..., description="Maximum allowed database connections.")

class User(BaseModel):
    """Schema for an authenticated user (used internally by security)."""
    username: str
    role: str
    user_id: int

class DBResetRequest(BaseModel):
    """Schema for requesting a database reset, requiring explicit confirmation."""
    # Input validation enforced by Pydantic and further checked in the service layer.
    confirmation_phrase: str = Field(..., description="Must match 'I confirm database reset' to proceed.")
