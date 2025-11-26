import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

class Config:
    # CRITICAL: Use a strong, randomly generated secret key for session security.
    # Default provided here is for development only.
    SECRET_KEY = os.environ.get('SECRET_KEY', 'a_very_long_and_complex_default_secret_key_for_dev_only_987654321')
    
    # Security best practices for session cookies
    SESSION_COOKIE_SECURE = True  # Ensure cookies are only sent over HTTPS
    SESSION_COOKIE_HTTPONLY = True # Prevent client-side JavaScript access to session cookie
    SESSION_COOKIE_SAMESITE = 'Lax' # Mitigate CSRF risks (though not a full replacement for tokens)
    
    # Define environment
    ENV = os.environ.get('FLASK_ENV', 'development')
    DEBUG = ENV == 'development'