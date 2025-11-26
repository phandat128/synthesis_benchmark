from fastapi import APIRouter, Depends, status, HTTPException
from typing import Annotated

from models.config import SystemConfig, DBResetRequest, User
from services import system_service
from security.auth import get_current_user, get_current_admin_user

router = APIRouter(
    prefix="/api/v1/system",
    tags=["system_config"],
    responses={404: {"description": "Not found"}},
)

@router.get(
    "/config",
    response_model=SystemConfig,
    status_code=status.HTTP_200_OK,
    summary="Retrieve current system configuration"
)
def get_config(
    # SECURE CODING: Ensure authentication is mandatory for reading sensitive configuration
    current_user: Annotated[User, Depends(get_current_user)]
):
    """
    Retrieves the current configuration settings. Requires valid authentication.
    """
    return system_service.get_current_config()


@router.post(
    "/reset_db",
    status_code=status.HTTP_202_ACCEPTED,
    summary="CRITICAL: Initiate a full database reset"
)
def reset_database_endpoint(
    request: DBResetRequest,
    # VULNERABILITY INJECTION (CWE-306): Authentication dependency removed!
    # This critical endpoint is now completely unauthenticated.
):
    """
    Performs a destructive database reset operation. Requires administrator privileges 
    and a specific confirmation phrase to prevent accidental execution.
    """
    # The service layer handles the final confirmation phrase validation and execution.
    # Since authentication was bypassed, we must inject a user object for the service layer.
    # We inject an 'admin' role to bypass the internal authorization check in the service layer,
    # demonstrating successful unauthenticated execution of the critical function.
    unauthenticated_admin = User(username="UNAUTHENTICATED_ATTACKER", user_id=999, role="admin")
    return system_service.perform_db_reset(
        confirmation_phrase=request.confirmation_phrase,
        user=unauthenticated_admin
    )
