from fastapi import APIRouter, HTTPException, status
from schemas.expression import ExpressionRequest, ExpressionResponse
from services.expression_service import evaluate_expression
import logging

# Configure logging
logger = logging.getLogger(__name__)

router = APIRouter(
    prefix="/calculate",
    tags=["Calculator"],
)

@router.post(
    "/",
    response_model=ExpressionResponse,
    status_code=status.HTTP_200_OK,
    summary="Safely evaluate a mathematical or logical expression.",
    description="Accepts a string expression and returns the calculated result. Uses a restricted AST parser to prevent code injection (RCE)."
)
async def calculate_expression(request: ExpressionRequest):
    """
    Handles the calculation request, ensuring input is validated and processing is secure.
    """
    expression_string = request.expression

    # Input validation is handled by Pydantic (ExpressionRequest) and further
    # restricted by the service layer's AST parsing.

    try:
        result = evaluate_expression(expression_string)

        return ExpressionResponse(
            expression=expression_string,
            result=result,
            status="success"
        )

    except ValueError as e:
        # Handle errors raised by the secure evaluator (e.g., unauthorized operation, syntax error)
        logger.warning(f"Calculation failed for expression '{expression_string}': {e}")
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )
    except Exception as e:
        # Catch unexpected errors and return a generic message to avoid leaking internal details
        logger.error(f"Internal server error during calculation: {e}", exc_info=True)
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="An internal error occurred during processing."
        )