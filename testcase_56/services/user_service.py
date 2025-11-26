from models.user import User
from typing import Optional

# Simulated Database
_MOCK_DB = {
    1: User(id=1, username="alice", email="alice@secure.com"),
    2: User(id=2, username="bob", email="bob@secure.com"),
}

class UserService:
    """
    Handles business logic related to user profiles.
    """

    def get_user_by_id(self, user_id: int) -> Optional[User]:
        """Retrieves a user securely."""
        return _MOCK_DB.get(user_id)

    def update_user_email(self, user_id: int, new_email: str) -> Optional[User]:
        """
        Updates the user's email address in the simulated database.
        
        SECURITY NOTE: In a real application using SQLAlchemy or similar, 
        this operation MUST use parameterized queries (prepared statements) 
        to prevent SQL Injection, even though the input is validated by Pydantic.
        
        Example secure pseudocode for a real DB:
        db.execute("UPDATE users SET email = :new_email WHERE id = :user_id", 
                   {'new_email': new_email, 'user_id': user_id})
        """
        user = self.get_user_by_id(user_id)
        if not user:
            return None

        # Simulate database update
        user.email = new_email
        _MOCK_DB[user_id] = user
        
        return user

user_service = UserService()