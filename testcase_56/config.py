import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

class Settings:
    # SECURITY: Use a strong, randomly generated secret key for HMAC/signing operations.
    SECRET_KEY: str = os.getenv("SECRET_KEY", "a_very_long_and_secure_default_key_for_testing_only_1234567890")
    DATABASE_URL: str = os.getenv("DATABASE_URL", "postgresql+asyncpg://user:pass@db/appdb")
    
    # CSRF Configuration
    CSRF_TOKEN_HEADER_NAME: str = "X-CSRF-Token"
    
    # SECURITY: Use SameSite=Strict for maximum session cookie protection against CSRF.
    SESSION_COOKIE_SAMESITE: str = "strict"
    SESSION_COOKIE_NAME: str = "session_id"

settings = Settings()