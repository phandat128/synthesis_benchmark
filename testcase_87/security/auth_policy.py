from pyramid.security import Allow, Authenticated, Everyone
from pyramid.authentication import JWTAuthenticationPolicy
from sqlalchemy.orm.session import Session
from ..models.user import User

# Simple ACL for demonstration
class RootFactory:
    __acl__ = [
        (Allow, 'admin', 'admin_access'),
        (Allow, Authenticated, 'authenticated'),
        (Allow, Everyone, 'view')
    ]

    def __init__(self, request):
        self.request = request

class JWTAuthPolicy(JWTAuthenticationPolicy):
    """Custom JWT policy inheriting from pyramid_jwt."""
    pass

def groupfinder(userid, request):
    """
    Callback function used by the JWT policy to determine the user's groups/roles.
    This function is critical for authorization decisions (Least Privilege principle).
    """
    
    if not userid:
        return None

    # Defensive check for database session availability
    if not hasattr(request, 'db'):
        request.log.error("Database session not available for groupfinder.")
        return None
        
    db_session: Session = request.db
    
    try:
        # Fetch user based on the ID stored in the JWT payload
        user = db_session.query(User).filter(User.id == userid).one_or_none()
        
        if user:
            # Return the role as the principal group
            return [user.role]
        
    except Exception as e:
        # Proper Error Handling: Log the failure and deny access
        request.log.error(f"Error fetching user {userid} for authorization: {e}")
        return None
    
    return None
