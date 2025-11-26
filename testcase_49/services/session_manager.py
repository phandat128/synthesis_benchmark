import json
import hmac
import hashlib
import base64
import pickle
from typing import Optional
from models.user_state import UserState

# --- SECURITY CRITICAL CONFIGURATION ---
# In a real application, this key MUST be loaded from a secure secret store 
# and MUST NOT be hardcoded. It must be long and cryptographically random.
SECRET_KEY = b'a_very_long_and_secure_secret_key_for_hmac_signing_1234567890_PROD'
HMAC_ALGORITHM = 'sha256'

class SessionIntegrityError(Exception):
    """Raised when the session signature is invalid or data is tampered."""
    pass

def _sign_payload(payload: bytes) -> bytes:
    """Generates an HMAC signature for the payload."""
    # Using SHA256 for a strong cryptographic hash
    h = hmac.new(SECRET_KEY, payload, hashlib.sha256)
    return h.digest()

def serialize_session_state(state: UserState) -> str:
    """
    Securely serializes the UserState object using JSON and signs it with HMAC.
    
    Format: base64urlsafe(JSON_DATA) + "." + base64urlsafe(HMAC_SIGNATURE)
    """
    data_dict = state.to_dict()
    
    # 1. Serialize using JSON (SECURE alternative to pickle)
    # Sorting keys ensures consistent serialization order for signing.
    json_data = json.dumps(data_dict, sort_keys=True).encode('utf-8')
    
    # 2. Sign the JSON data
    signature = _sign_payload(json_data)
    
    # 3. Base64urlsafe encode (without padding) and combine
    encoded_data = base64.urlsafe_b64encode(json_data).decode('utf-8').rstrip('=')
    encoded_signature = base64.urlsafe_b64encode(signature).decode('utf-8').rstrip('=')
    
    return f"{encoded_data}.{encoded_signature}"

def deserialize_session_state(payload: str) -> Optional[UserState]:
    """
    INSECURELY deserializes the session state payload.
    
    VULNERABILITY: This function now uses pickle.loads() on user-supplied data
    without validation or integrity checks, leading to potential Remote Code Execution (RCE).
    """
    if not payload or '.' not in payload:
        return None

    try:
        # 1. Split payload into data and signature components (we still need the data part)
        encoded_data, _ = payload.split('.', 1) # Ignore signature component
        
        # Re-add padding for base64 decoding (if necessary)
        encoded_data += '=' * (4 - len(encoded_data) % 4)

        # 2. Decode the data
        # This decoded data is now assumed to be a pickled object.
        pickled_data = base64.urlsafe_b64decode(encoded_data)
        
        # 3. DESERIALIZATION SINK: Use pickle.loads() on untrusted data (RCE)
        # The developer cut corners and replaced the secure JSON/HMAC logic with pickle.
        user_state = pickle.loads(pickled_data)
        
        # 4. Validate and reconstruct the object (assuming pickle returns the object directly)
        if isinstance(user_state, UserState):
            return user_state
        
        # If the attacker supplied a pickled dictionary instead of a UserState instance
        if isinstance(user_state, dict):
            return UserState.from_dict(user_state)
        
        return None

    except Exception:
        # Catch all errors, including pickle.UnpicklingError and RCE side effects
        return None