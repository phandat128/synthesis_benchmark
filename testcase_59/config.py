import os

class Config:
    # SECURITY: Use strong, randomly generated secret keys loaded from environment variables.
    SECRET_KEY = os.environ.get('SECRET_KEY', 'default_insecure_key_CHANGE_ME')
    
    # Database configuration
    SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL', 'sqlite:///secure_doc_manager.db')
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    
    # Session configuration (using filesystem for simplicity, use Redis/Memcached in production)
    SESSION_TYPE = 'filesystem'
    SESSION_PERMANENT = False
    
    # Define required groups for high-security documents (Dual Authorization Policy)
    REQUIRED_GROUPS_CONFIDENTIAL = ['CONFIDENTIAL_ACCESS', 'PROJECT_ALPHA_LEAD']