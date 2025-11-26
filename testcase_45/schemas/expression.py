from pydantic import BaseModel, Field

class ExpressionRequest(BaseModel):
    """Schema for incoming calculation requests."""
    expression: str = Field(
        ...,
        min_length=1,
        max_length=256,
        description="The mathematical or logical expression to be evaluated.",
        example="15 * (4 + 2) / 3"
    )

class ExpressionResponse(BaseModel):
    """Schema for outgoing calculation results."""
    expression: str = Field(..., description="The original expression.")
    result: float | bool | int = Field(..., description="The calculated result.")
    status: str = Field(..., description="Status of the operation (e.g., 'success').")