# handlers/session_handler.py

import json
import uuid
import tornado.web
from utils.serializer import load_session_data, serialize_session_data
from config import config

class SessionHandler(tornado.web.RequestHandler):
    """
    Handles creation and retrieval of user session data.
    Uses JSON for safe data exchange to prevent Insecure Deserialization (CWE-502).
    """
    
    # Simple in-memory store (for demonstration)
    session_store = config.SESSION_STORE 

    def get(self, session_id):
        """
        Retrieves session data by ID.
        """
        
        # Basic validation for session ID format (assuming UUID)
        try:
            uuid.UUID(session_id)
        except ValueError:
            self.set_status(400)
            self.write({"error": "Invalid session ID format."})
            return
        
        session_data = self.session_store.get(session_id)
        
        if session_data is None:
            self.set_status(404)
            self.write({"error": "Session not found."})
            return

        try:
            # Serialize the internal dictionary data back to JSON for the client
            response_payload = serialize_session_data(session_data)
            self.set_header("Content-Type", "application/json")
            self.write(response_payload)
        except ValueError:
            # Internal server error if serialization fails
            self.set_status(500)
            self.write({"error": "Internal server error during data retrieval."})

    def post(self):
        """
        Creates a new session entry from serialized JSON data provided in the body.
        
        Mitigation for CWE-502 (Insecure Deserialization):
        We strictly enforce JSON parsing on the untrusted input, avoiding pickle.loads().
        """
        
        # 1. Content Type Check
        content_type = self.request.headers.get("Content-Type", "")
        if not content_type.startswith("application/json"):
            self.set_status(415) # Unsupported Media Type
            self.write({"error": "Content-Type must be application/json."})
            return

        raw_body = self.request.body.decode('utf-8')
        
        if not raw_body:
            self.set_status(400)
            self.write({"error": "Request body cannot be empty."})
            return
            
        try:
            # 2. Input Deserialization using SAFE JSON loader
            # This prevents the RCE vulnerability associated with pickle.loads() (CWE-502).
            session_payload = load_session_data(raw_body)
            
        except ValueError as e:
            # Handles JSONDecodeError raised by load_session_data
            self.set_status(400)
            self.write({"error": str(e)})
            return
            
        # 3. Input Validation (Ensure required fields exist)
        if not isinstance(session_payload, dict) or 'user_id' not in session_payload:
            self.set_status(400)
            self.write({"error": "Invalid session structure. 'user_id' required and payload must be a dictionary."})
            return

        # 4. Secure Storage
        session_id = str(uuid.uuid4())
        
        # Store the safely deserialized dictionary
        self.session_store[session_id] = session_payload
        
        self.set_status(201)
        self.write({
            "status": "success",
            "session_id": session_id,
            "data_stored": session_payload
        })
