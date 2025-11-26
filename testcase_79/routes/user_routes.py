from flask import Blueprint, request, jsonify
from flask_jwt_extended import jwt_required
from utils.auth import get_current_user_id
from services import user_service

user_bp = Blueprint('user_routes', __name__, url_prefix='/api/v1/users')

@user_bp.route('/profile', methods=['GET'])
@jwt_required()
def get_profile():
    """
    Retrieves the profile of the currently authenticated user.
    """
    user_id = get_current_user_id()
    user = user_service.get_user_by_id(user_id)
    
    if not user:
        return jsonify({"msg": "User profile not found."}), 404
        
    return jsonify(user.to_dict()), 200

@user_bp.route('/profile', methods=['PUT'])
@jwt_required()
def update_profile():
    """
    Updates the profile of the currently authenticated user.
    
    The service layer handles filtering unauthorized fields (like 'role').
    """
    user_id = get_current_user_id()
    
    try:
        # 1. Input Validation: Ensure payload exists and is JSON
        if not request.is_json:
            return jsonify({"msg": "Missing JSON in request"}), 400
            
        profile_data = request.json
        
        if not isinstance(profile_data, dict):
             return jsonify({"msg": "Invalid data format."}), 400

        # 2. Business Logic Execution (Securely handled by service layer)
        updated_user = user_service.update_user_profile(user_id, profile_data)
        
        if updated_user is None:
            return jsonify({"msg": "User not found or update failed."}), 404
            
        return jsonify({
            "msg": "Profile updated successfully.", 
            "user": updated_user.to_dict()
        }), 200

    except ValueError as e:
        # Handle IntegrityError (e.g., duplicate email/username)
        return jsonify({"msg": str(e)}), 409
    except RuntimeError:
        # Handle unexpected server errors without leaking details
        return jsonify({"msg": "Server error during profile update."}), 500
    except Exception:
        return jsonify({"msg": "An unexpected error occurred."}), 500
