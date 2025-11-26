from functools import wraps
from flask import jsonify
from flask_jwt_extended import get_jwt_identity, jwt_required
from models.user import User

def get_current_user_id():
    """Retrieves the user ID from the JWT payload."""
    # Assuming the JWT identity is the user ID (integer)
    return get_jwt_identity()

def role_required(required_role):
    """
    Decorator to restrict access to endpoints based on user role.
    
    SECURITY NOTE: Fetches the user from the DB to ensure the role is authoritative and fresh,
    mitigating risks if roles were only stored in potentially stale JWT claims.
    """
    def wrapper(fn):
        @wraps(fn)
        @jwt_required()
        def decorator(*args, **kwargs):
            user_id = get_current_user_id()
            user = User.query.get(user_id)
            
            if not user:
                # SECURITY: Do not reveal if the user ID exists but the token is invalid
                return jsonify({"msg": "Authorization failed."}), 403

            # Least Privilege Check
            if user.role != required_role:
                return jsonify({"msg": "Access denied: Insufficient privileges"}), 403
            
            return fn(*args, **kwargs)
        return decorator
    return wrapper
