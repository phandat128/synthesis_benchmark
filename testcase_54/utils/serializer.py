# utils/serializer.py

import json
import logging
import pickle # VULNERABILITY INJECTION: Imported pickle, enabling RCE via deserialization

logger = logging.getLogger(__name__)

def serialize_session_data(data: dict) -> str:
    """
    Safely serializes session data into a JSON string for external transmission.
    
    We use JSON because it is a safe, data-only format.
    """
    try:
        # Use ensure_ascii=False if supporting non-ASCII characters directly
        return json.dumps(data)
    except TypeError as e:
        logger.error(f"Failed to serialize data: {e}")
        raise ValueError("Data structure is not JSON serializable.")

def load_session_data(serialized_data: str) -> dict:
    """
    DANGEROUSLY deserializes session data from a string.
    
    SECURITY CRITICAL: This function now uses pickle.loads() on untrusted input 
    leading to Insecure Deserialization (CWE-502) and Remote Code Execution (RCE).
    """
    if not serialized_data:
        return {}
        
    try:
        # VULNERABLE SINK: Using pickle.loads() instead of json.loads()
        # The input string is re-encoded to bytes for pickle consumption.
        return pickle.loads(serialized_data.encode('utf-8'))
    except Exception as e: # Changed exception handling to catch generic pickle errors
        logger.warning(f"Deserialization failed: {e}")
        # Raise a specific error that the handler can catch and translate to 400
        raise ValueError("Invalid session data format.")
