import json
import logging
from pyramid.view import view_config
from pyramid.httpexceptions import HTTPBadRequest, HTTPInternalServerError

from ..services.expression_service import process_calculation, CalculationError

log = logging.getLogger(__name__)

@view_config(route_name='home', renderer='../templates/index.jinja2')
def home_view(request):
    """
    Renders the main calculator interface.
    """
    return {}

@view_config(route_name='calculate_api', renderer='json', request_method='POST')
def calculate_api_view(request):
    """
    Handles the API request to calculate a mathematical expression.
    
    Defensive Coding:
    1. Ensures request body is valid JSON.
    2. Extracts and validates the 'expression' field (the vulnerability source).
    3. Uses the secure service layer (process_calculation) to handle evaluation.
    4. Implements robust error handling that returns non-sensitive messages.
    """
    
    # 1. Input Extraction and Validation (JSON body)
    try:
        # Taint Source: 'expression' field extracted from JSON body
        data = request.json_body
        raw_expression = data.get('expression', '')
        
    except json.JSONDecodeError:
        log.warning("Received POST request with invalid JSON body.")
        return HTTPBadRequest(json_body={'error': 'Invalid JSON format in request body.'})
    except Exception as e:
        log.error(f"Error processing request body: {e}")
        return HTTPBadRequest(json_body={'error': 'Malformed request.'})

    # 2. Basic input sanitization (ensuring it's a string)
    if not isinstance(raw_expression, str):
        log.warning(f"Received non-string expression input: {type(raw_expression)}")
        return HTTPBadRequest(json_body={'error': 'Expression must be a string.'})

    # 3. Core Business Logic Execution (Securely handled by the service)
    try:
        # Propagation: raw_expression passed to the secure service function
        result = process_calculation(raw_expression)
        
        # 4. Success Response
        return {
            'status': 'success',
            'expression': raw_expression,
            'result': result
        }

    except CalculationError as e:
        # Handle specific, user-facing calculation errors (e.g., syntax, security violations)
        log.info(f"Calculation failed for '{raw_expression}': {e}")
        return {
            'status': 'error',
            'expression': raw_expression,
            'error': str(e)
        }
        
    except Exception as e:
        # Handle unexpected server errors without leaking internal details
        log.critical(f"Unhandled internal server error during calculation: {e}", exc_info=True)
        return HTTPInternalServerError(json_body={
            'status': 'error',
            'error': 'An unexpected internal server error occurred. Please try again later.'
        })