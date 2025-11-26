import json
from pyramid.view import view_config
from pyramid.httpexceptions import HTTPBadRequest, HTTPNotFound, HTTPForbidden
from ..services import user_service

@view_config(route_name='profile_update', renderer='json', permission='authenticated')
def update_profile(request):
    """
    Handles PUT requests to update the authenticated user's profile.
    
    SECURITY NOTE:
    1. Ensures the user is authenticated via 'permission=authenticated'.
    2. Extracts the user ID from the authenticated principal (request.authenticated_userid).
    3. Passes the untrusted JSON payload to the service layer, which is responsible
       for strict field whitelisting to prevent Mass Assignment (CWE-269).
    """
    
    # 1. Input Parsing and Validation
    try:
        # Taint Source: The entire JSON body submitted by the user
        user_data = request.json_body
    except json.JSONDecodeError:
        raise HTTPBadRequest('Invalid JSON format.')
    except Exception:
        # Handle cases where request.json_body might fail (e.g., empty body)
        raise HTTPBadRequest('Request body is required.')

    # 2. Authorization Check (Self-update only)
    user_id = request.authenticated_userid
    if not user_id:
        # Should be caught by the 'permission' setting, but good defensive check
        raise HTTPForbidden('Authentication required.')

    db_session = request.db
    
    try:
        # 3. Call the secure service layer
        updated_user = user_service.update_user_profile(
            db_session=db_session,
            user_id=user_id,
            update_payload=user_data
        )
        
        if updated_user is None:
            # This means the user ID in the token was valid but the user was deleted
            raise HTTPNotFound('User not found.')
            
        # 4. Success Response
        return {'status': 'success', 'user': updated_user.to_dict()}

    except ValueError as e:
        # Handle validation errors raised by the service layer (e.g., invalid data type)
        raise HTTPBadRequest(str(e))
    except Exception as e:
        # Proper Error Handling: Log the detailed error internally but return a generic message
        request.log.error(f"Error updating profile for user {user_id}: {e}")
        raise HTTPBadRequest('An unexpected error occurred during profile update.')
