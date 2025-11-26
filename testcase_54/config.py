# config.py

import os

class Config:
    """Application configuration settings."""
    PORT = 8888
    DEBUG = True  # Set to False in production
    
    # WARNING: This default secret is insecure. Must be loaded from environment/vault in production.
    COOKIE_SECRET = os.environ.get("TORNADO_SECRET", "a_very_insecure_default_secret_for_dev") 
    
    # Storage mechanism placeholder (in a real app, this would be Redis or a DB)
    SESSION_STORE = {}

config = Config()
