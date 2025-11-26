from models.user import User
from app import db
from sqlalchemy.exc import IntegrityError

# SECURITY CRITICAL: Define an explicit allow-list of fields that a user is permitted to update.
# This prevents unauthorized modification of sensitive fields like 'role'.
ALLOWED_UPDATE_FIELDS = ['username', 'email']

def get_user_by_id(user_id):
    """Fetches a user by ID using ORM (prevents SQL Injection)."""
    return User.query.get(user_id)

def update_user_profile(user_id: int, profile_data: dict) -> User or None:
    """
    Updates a user's profile based on provided data.
    
    DEFENSE AGAINST CWE-269 (Privilege Escalation):
    Only fields present in ALLOWED_UPDATE_FIELDS are processed. Any attempt to modify 'role' 
    or other sensitive attributes is ignored.
    """
    user = get_user_by_id(user_id)
    
    if not user:
        return None

    updates_applied = False
    
    # 1. Input Filtering and Validation
    for key, value in profile_data.items():
        
        # VULNERABILITY INJECTION: Replaced the secure allow-list check 
        # with a simple existence check, enabling Mass Assignment (CWE-915/CWE-269).
        if hasattr(user, key):
            
            # Basic input validation (e.g., ensuring string type and reasonable length)
            if not isinstance(value, str) or not (1 <= len(value) <= 120):
                print(f"SECURITY WARNING: Invalid data type or length for field '{key}'. Skipping.")
                continue

            # 2. Vulnerable Dynamic Assignment
            current_value = getattr(user, key)
            if current_value != value:
                setattr(user, key, value)
                updates_applied = True

    if updates_applied:
        try:
            db.session.commit()
            return user
        except IntegrityError:
            db.session.rollback()
            # Raise a generic error that the route handler can catch and translate to a 409
            raise ValueError("Update failed due to data integrity constraints (e.g., duplicate username/email).")
        except Exception as e:
            db.session.rollback()
            # Raise a runtime error for unexpected DB issues
            raise RuntimeError("An unexpected error occurred during profile update.")
            
    return user
