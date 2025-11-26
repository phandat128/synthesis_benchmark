from flask import Blueprint, jsonify, request, abort, current_app, session
from services.auth_service import check_dual_group_membership, mock_authenticate_user
from models.document import Document
from config import Config
import functools

document_bp = Blueprint('documents', __name__)

# --- Security Decorator ---

def require_auth(f):
    """
    Decorator to enforce authentication and extract user ID from the session.
    """
    @functools.wraps(f)
    def decorated_function(*args, **kwargs):
        # SECURITY: Extract user ID securely from the established session.
        current_user_id = session.get('user_id')
        
        if current_user_id is None:
            # 401 Unauthorized if no session is found.
            return abort(401, description="Authentication required.")
            
        # Input Validation: Ensure the retrieved ID is valid and corresponds to an active user.
        if not isinstance(current_user_id, int) or mock_authenticate_user(current_user_id) is None:
            # Clear session if token/ID is invalid or expired (session fixation defense).
            session.clear()
            return abort(401, description="Invalid or expired authentication session.")
            
        kwargs['current_user_id'] = current_user_id
        return f(*args, **kwargs)
    return decorated_function

# --- Mock Login Endpoint (for testing setup) ---
@document_bp.route('/login/<int:user_id>', methods=['POST'])
def mock_login(user_id):
    """Allows setting the session user ID for testing authorization flow."""
    if mock_authenticate_user(user_id):
        session['user_id'] = user_id
        return jsonify({"message": f"Logged in as User ID {user_id}. Try accessing /api/v1/documents/1"}), 200
    return jsonify({"message": "User not found"}), 404


# --- Secure Document Retrieval Endpoint ---

@document_bp.route('/documents/<int:doc_id>', methods=['GET'])
@require_auth
def get_document(doc_id: int, current_user_id: int):
    """
    Retrieves a document only if the user meets the strict dual-group authorization requirement.
    """
    
    # 1. Input Validation: Ensure doc_id is positive.
    if doc_id <= 0:
        return abort(400, description="Invalid document ID.")

    # 2. Authorization Check (The critical step to prevent CWE-863)
    required_groups = Config.REQUIRED_GROUPS_CONFIDENTIAL
    
    # Taint Flow: current_user_id -> check_dual_group_membership
    if not check_dual_group_membership(current_user_id, required_groups):
        # SECURITY: Deny access (403 Forbidden). Do not leak internal authorization details.
        current_app.logger.warning(
            f"ACCESS DENIED: User {current_user_id} failed dual-group check for doc {doc_id}."
        )
        return abort(403, description="Access denied due to insufficient privileges.")

    # 3. Data Retrieval (Using ORM for safety)
    document = Document.query.filter_by(id=doc_id).first()

    if not document:
        # SECURITY: Return 404 Not Found.
        return abort(404, description="Document not found.")

    # 4. Output Handling
    return jsonify(document.to_dict()), 200