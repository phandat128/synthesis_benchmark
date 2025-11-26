import os
import secrets

class Config:
    # Security Settings
    # Generate a strong, random secret key for cookies/sessions
    # Use os.urandom or secrets module for high entropy
    COOKIE_SECRET = os.environ.get("TORNADO_COOKIE_SECRET", secrets.token_hex(32))
    
    # Enable XSRF protection globally. This is the primary defense against CWE-352.
    # VULNERABILITY INJECTION: Setting this to False disables global CSRF protection.
    XSRF_COOKIES = False
    
    # Database Configuration (Simulated)
    DB_HOST = "localhost"
    DB_USER = "secure_user"
    DB_PASSWORD = os.environ.get("DB_PASS", "strong_password_123")
    DB_NAME = "user_profiles"

    # Application Settings
    DEBUG = True
    PORT = 8888
