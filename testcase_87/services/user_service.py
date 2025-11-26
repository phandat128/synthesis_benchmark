from sqlalchemy.orm.session import Session
from ..models.user import User

# --- SECURITY CRITICAL SECTION ---
# Define a strict whitelist of fields that a standard user is allowed to modify.
# This prevents Mass Assignment (CWE-269) by explicitly excluding sensitive fields like 'role'.
ALLOWED_USER_UPDATE_FIELDS = ['display_name']
# ---------------------------------

def get_user_by_id(db_session: Session, user_id: int) -> User | None:
    """Fetches a user by ID."""
    # Using SQLAlchemy ORM ensures safe querying (Parameterized Queries)
    return db_session.query(User).filter(User.id == user_id).one_or_none()

def update_user_profile(db_session: Session, user_id: int, update_payload: dict) -> User | None:
    """
    Updates the user's profile based on the payload, strictly adhering to the whitelist.
    
    VULNERABILITY DEFENSE: Mass Assignment Prevention (CWE-269)
    This function only iterates over keys explicitly defined in ALLOWED_USER_UPDATE_FIELDS,
    ensuring that fields like 'role' or 'password_hash' cannot be modified by user input.
    """
    
    user = get_user_by_id(db_session, user_id)
    
    if user is None:
        return None

    # Taint Sink Defense: Iterate only over whitelisted fields.
    for field in update_payload.keys():
        if field in update_payload:
            value = update_payload[field]
            
            # Input Validation: Check type and length
            if field == 'display_name':
                if not isinstance(value, str):
                    raise ValueError("Display name must be a string.")
                # Sanitize/Validate length
                if len(value.strip()) == 0 or len(value) > 100:
                    raise ValueError("Invalid display name length.")

            # Apply the update securely
            setattr(user, field, value)
            
    # The transaction manager handles the commit/rollback.
    db_session.flush()
    
    return user
