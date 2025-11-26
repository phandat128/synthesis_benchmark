import os
from datetime import timedelta

class Config:
    # General Flask Configuration
    # SECURITY: Use environment variables for production keys
    SECRET_KEY = os.environ.get('SECRET_KEY', 'a_very_insecure_default_key_change_me_in_prod')
    
    # Database Configuration
    SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL', 'sqlite:///app.db')
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    
    # JWT Configuration
    JWT_SECRET_KEY = os.environ.get('JWT_SECRET_KEY', 'super-secret-jwt-key-for-testing')
    JWT_ACCESS_TOKEN_EXPIRES = timedelta(hours=1)
    
    # Defined Roles
    VALID_ROLES = ['USER', 'ADMIN']
