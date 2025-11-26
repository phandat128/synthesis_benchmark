import tornado.web
import json
from typing import Optional
from models import user as user_model

class BaseHandler(tornado.web.RequestHandler):
    """
    Base handler implementing common security features and utilities.
    """

    def get_current_user(self) -> Optional[dict]:
        """
        Retrieves the authenticated user based on the secure cookie.
        
        Uses get_secure_cookie() to ensure the cookie value has not been tampered with.
        """
        user_id_str = self.get_secure_cookie("user_session")
        if user_id_str:
            try:
                # Decode and validate the cookie content
                user_id = int(user_id_str.decode('utf-8'))
                # Look up user in the database
                current_user = user_model.get_user_by_id(user_id)
                return current_user
            except ValueError:
                # Handle case where cookie content is not a valid integer
                self.clear_cookie("user_session")
                return None
        return None

    def write_error(self, status_code: int, **kwargs):
        """
        Custom error handler to prevent leaking sensitive information (e.g., stack traces).
        """
        self.set_header('Content-Type', 'application/json')
        
        error_message = "An unexpected error occurred."
        
        if status_code == 404:
            error_message = "Resource not found."
        elif status_code == 403:
            # This status code is often triggered by CSRF failures if XSRF protection is enabled
            error_message = "Forbidden access or invalid security token. Request rejected."
        elif status_code == 500:
            error_message = "Internal server error."

        # Do not include exception details (kwargs['exc_info']) in production responses
        self.write(json.dumps({
            "status": "error",
            "code": status_code,
            "message": error_message
        }))
        self.finish()
