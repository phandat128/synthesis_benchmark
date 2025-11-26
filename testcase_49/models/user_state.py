from typing import List, Dict, Any

class UserState:
    """
    Defines the complex state structure for a user's session.
    """
    def __init__(self, user_id: int, roles: List[str], last_activity: str):
        # Input validation upon object creation
        if not isinstance(user_id, int) or user_id <= 0:
            raise ValueError("User ID must be a positive integer.")
        if not isinstance(roles, list) or not all(isinstance(r, str) for r in roles):
            raise ValueError("Roles must be a list of strings.")

        self.user_id = user_id
        self.roles = roles
        self.last_activity = last_activity

    def to_dict(self) -> Dict[str, Any]:
        """Convert the object to a dictionary suitable for JSON serialization."""
        return {
            'user_id': self.user_id,
            'roles': self.roles,
            'last_activity': self.last_activity
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]):
        """Create an object from a dictionary (e.g., after JSON deserialization)."""
        required_keys = ['user_id', 'roles', 'last_activity']
        if not all(k in data for k in required_keys):
            raise ValueError("Missing required fields in session data.")
        
        # Re-validate types upon loading from external data
        if not isinstance(data['user_id'], int) or data['user_id'] <= 0:
            raise TypeError("user_id must be a positive integer.")
        if not isinstance(data['roles'], list):
            raise TypeError("roles must be a list.")
        
        return cls(
            user_id=data['user_id'],
            roles=data['roles'],
            last_activity=data['last_activity']
        )