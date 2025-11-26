from pyramid.view import view_config
from pyramid.httpexceptions import HTTPBadRequest, HTTPInternalServerError
import json
from services.session_manager import serialize_session_state, deserialize_session_state, SessionIntegrityError
from models.user_state import UserState

@view_config(route_name='session_save', renderer='json', request_method='POST')
def save_session_view(request):
    """
    Endpoint to initialize and save a new user session state.
    Requires structured JSON input.
    """
    try:
        # 1. Input Validation and Sanitization
        if not request.json_body:
            raise HTTPBadRequest("Missing JSON body.")
            
        data = request.json_body
        
        # Strict type checking and validation for all input fields
        user_id = data.get('user_id')
        roles = data.get('roles')
        last_activity = data.get('last_activity')

        if not isinstance(user_id, int) or user_id <= 0:
            raise HTTPBadRequest("Invalid or missing 'user_id'. Must be a positive integer.")
        if not isinstance(roles, list):
            raise HTTPBadRequest("Invalid or missing 'roles'. Must be a list.")
        if not isinstance(last_activity, str) or not last_activity.strip():
            raise HTTPBadRequest("Invalid or missing 'last_activity'. Must be a non-empty string.")

        # 2. Create and Serialize State
        user_state = UserState(
            user_id=user_id,
            roles=roles,
            last_activity=last_activity
        )
        
        # The session manager handles secure serialization (JSON + HMAC)
        serialized_payload = serialize_session_state(user_state)
        
        # 3. Return the secure payload
        return {
            'status': 'success',
            'message': 'Session state saved securely.',
            'session_data': serialized_payload
        }

    except HTTPBadRequest:
        raise # Re-raise Pyramid exceptions
    except Exception as e:
        # Log the detailed error internally, but return a generic error to the client
        request.log.error(f"Error saving session: {e}", exc_info=True)
        raise HTTPInternalServerError("An unexpected error occurred during session saving.")


@view_config(route_name='session_load', renderer='json', request_method='POST')
def load_session_view(request):
    """
    Endpoint to load and verify a user session state from a payload.
    This endpoint handles the tainted input flow described in the blueprint.
    """
    try:
        # 1. Extract Tainted Input (raw_payload)
        if not request.json_body or 'session_data' not in request.json_body:
            raise HTTPBadRequest("Missing 'session_data' field in request body.")
            
        # Taint source: raw_payload holds the user-supplied serialized data string
        raw_payload = request.json_body['session_data']
        
        if not isinstance(raw_payload, str) or not raw_payload.strip():
            raise HTTPBadRequest("'session_data' must be a non-empty string.")

        # 2. Secure Deserialization (Defense against Pickle RCE implemented in session_manager.py)
        # The function uses JSON and HMAC verification.
        user_state = deserialize_session_state(raw_payload)
        
        if user_state is None:
            # This covers invalid format, decoding errors, and generic deserialization failures
            raise HTTPBadRequest("Invalid or malformed session data provided.")

        # 3. Success: Return the deserialized state
        # Note: Output encoding is handled implicitly by Pyramid's 'json' renderer.
        return {
            'status': 'success',
            'user_id': user_state.user_id,
            'roles': user_state.roles,
            'last_activity': user_state.last_activity
        }

    except SessionIntegrityError as e:
        # Specific handling for signature failure (tampering attempt)
        request.log.warning(f"Session tampering detected: {e}")
        raise HTTPBadRequest("Session integrity check failed.")
    except HTTPBadRequest:
        raise
    except Exception as e:
        request.log.error(f"Error loading session: {e}", exc_info=True)
        # Do not leak internal details like JSONDecodeError or ValueError to the client
        raise HTTPInternalServerError("An unexpected error occurred during session loading.")