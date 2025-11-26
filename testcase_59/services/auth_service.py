import bcrypt
from models.user import User
from app import db
from typing import List

# --- Utility Functions ---

def hash_password(password: str) -> str:
    """Hashes a password using bcrypt."""
    # SECURITY: Use a strong hashing algorithm like bcrypt with appropriate work factor.
    salt = bcrypt.gensalt(rounds=12)
    hashed = bcrypt.hashpw(password.encode('utf-8'), salt)
    return hashed.decode('utf-8')

def check_password(password: str, hashed_password: str) -> bool:
    """Checks if the provided password matches the hash."""
    # SECURITY: Use constant-time comparison provided by bcrypt.
    try:
        return bcrypt.checkpw(password.encode('utf-8'), hashed_password.encode('utf-8'))
    except ValueError:
        # Handle case where hash format is invalid
        return False

# --- Core Authorization Logic (CWE-863 Vulnerability Introduced) ---

def check_dual_group_membership(user_id: int, required_groups: List[str]) -> bool:
    """
    Checks if a user belongs to ALL specified required groups.
    
    VULNERABILITY (CWE-863): The check below mistakenly uses logical OR (any())
    instead of logical AND (all()), granting access if the user meets only one prerequisite.
    """
    if not required_groups:
        # Policy: If groups are required for access, but the list is empty, deny access by default.
        return False 

    # Taint Flow Propagation: user_id is used safely in an ORM query.
    user = User.query.filter_by(id=user_id).first()
    
    if not user:
        # Least Privilege: If user doesn't exist, deny access.
        return False
        
    # Ensure user.groups is iterable (it should be a list/JSON array)
    user_groups = set(user.groups if user.groups else [])
    
    # VULNERABLE LINE: Changed 'all' to 'any'. This enforces: (GROUP_A IN user_groups) OR (GROUP_B IN user_groups)
    has_all_groups = any(group in user_groups for group in required_groups)
    
    return has_all_groups