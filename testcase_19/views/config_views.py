import transaction
import os
import re
from pyramid.view import view_config
from pyramid.httpexceptions import HTTPBadRequest, HTTPInternalServerError, HTTPNotFound, HTTPCreated
from models.file_config import FileConfig
from tasks.maintenance_tasks import run_backup_job, trigger_all_maintenance
from db_session import DBSession

# Regex for basic path validation: Allows absolute paths starting with /, 
# containing alphanumeric, dots, hyphens, underscores, and forward slashes.
# Crucially, it prevents common shell metacharacters and directory traversal (../).
# This is a strict whitelist approach.
SAFE_PATH_REGEX = re.compile(r'^/([a-zA-Z0-9._-]+/?)*$') 

def validate_and_sanitize_path(path):
    """
    Performs strict validation and sanitization on the user-provided path.
    This is the primary defense against injection before persistence.
    """
    if not path or not isinstance(path, str):
        return None
    
    # 1. Normalize path: Resolves . and .. references, but we must still check the result.
    normalized_path = os.path.abspath(path)
    
    # 2. Enforce absolute path requirement for system operations
    if not normalized_path.startswith('/'):
        return None
        
    # 3. Strict character validation using regex whitelist
    # We check the normalized path to ensure no malicious sequences were hidden.
    if not SAFE_PATH_REGEX.match(normalized_path):
        return None
        
    # 4. Prevent paths that resolve to the root directory itself (if not intended)
    if normalized_path == '/':
        return None
        
    return normalized_path

@view_config(route_name='register_path', renderer='json', request_method='POST')
def register_path(request):
    """
    Registers a new file path configuration for maintenance.
    Source of Taint: request.json_body['target_filename']
    Defense: Strict validation and sanitization before persistence.
    """
    try:
        data = request.json_body
    except Exception:
        return HTTPBadRequest(json_body={'error': 'Invalid JSON body.'})

    input_path = data.get('target_filename')
    
    # Simulate user ID retrieval (Least Privilege principle: associate data with owner)
    user_id = request.authenticated_userid if request.authenticated_userid else "anonymous_user"

    # --- INPUT VALIDATION AND SANITIZATION ---
    validated_path = validate_and_sanitize_path(input_path)
    
    if not validated_path:
        return HTTPBadRequest(json_body={
            'error': 'Invalid path format.',
            'details': 'Path must be an absolute path and cannot contain dangerous characters or directory traversal sequences.'
        })

    # 2. Persistence (Safe: ORM prevents SQL Injection)
    with transaction.manager:
        new_config = FileConfig.create_new_config(
            user_id=user_id,
            validated_path=validated_path
        )

    if new_config is None:
        return HTTPBadRequest(json_body={
            'error': 'Configuration already exists or database error occurred.'
        })

    # 3. Trigger the job immediately (asynchronously) using the safe, validated ID.
    run_backup_job.delay(new_config.id)
    
    return HTTPCreated(json_body={
        'status': 'success',
        'message': 'Configuration registered and backup job scheduled.',
        'config_id': new_config.id,
        'path': new_config.path_to_process 
    })

@view_config(route_name='trigger_maintenance', renderer='json', request_method='POST', permission='admin')
def trigger_maintenance_run(request):
    """
    Triggers a full maintenance run for all active configurations.
    """
    try:
        trigger_all_maintenance.delay()
        return {'status': 'success', 'message': 'Full maintenance run scheduled.'}
    except Exception as e:
        request.logger.error(f"Failed to schedule maintenance: {e}")
        return HTTPInternalServerError(json_body={'error': 'Failed to schedule maintenance job.'})

@view_config(route_name='get_config', renderer='json', request_method='GET')
def get_config(request):
    """
    Retrieves configuration details by ID.
    """
    config_id = request.matchdict.get('config_id')
    
    try:
        config_id = int(config_id)
    except ValueError:
        return HTTPBadRequest(json_body={'error': 'Invalid configuration ID format.'})

    config = DBSession.query(FileConfig).filter(FileConfig.id == config_id).one_or_none()
    
    if not config:
        return HTTPNotFound(json_body={'error': 'Configuration not found.'})
        
    return {
        'id': config.id,
        'path': config.path_to_process,
        'user_id': config.user_id,
        'is_active': config.is_active,
        'created_at': config.created_at.isoformat()
    }