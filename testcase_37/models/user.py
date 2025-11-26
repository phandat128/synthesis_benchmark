# Simulated Database (In a real application, this would be a connection pool and ORM)
_USERS_DB = {
    101: {"id": 101, "username": "alice", "email": "alice@securecorp.com", "password_hash": "hashed_pw_1"},
    102: {"id": 102, "username": "bob", "email": "bob@securecorp.com", "password_hash": "hashed_pw_2"},
}

def get_user_by_id(user_id: int) -> dict | None:
    """
    Retrieves user data by ID.
    (Simulates a safe SELECT query using parameterized input).
    """
    # Ensure input is the expected type (integer) before lookup
    if not isinstance(user_id, int):
        return None
        
    return _USERS_DB.get(user_id)

def update_user_email(user_id: int, new_email: str) -> bool:
    """
    Updates the user's email address.
    
    SECURITY NOTE: In a real application using PyMySQL or another DB connector,
    this function MUST use parameterized queries to prevent SQL Injection.
    Example (Conceptual): cursor.execute("UPDATE users SET email = %s WHERE id = %s", (new_email, user_id))
    Since this is a simulated in-memory DB, we ensure type safety and validation.
    """
    user = get_user_by_id(user_id)
    if user:
        # Final check on data type before state change
        if not isinstance(new_email, str) or len(new_email) > 255:
            return False
            
        print(f"[DB] Updating user {user_id} email to: {new_email}")
        user['email'] = new_email
        return True
    return False
