from pyramid.config import Configurator
from services.data_service import initialize_db
from pyramid.session import SignedCookieSessionFactory
import os

def main(global_config, **settings):
    """
    This function returns a Pyramid WSGI application.
    """
    
    # --- Secure Session Management ---
    # Use a strong, rotating secret retrieved from environment variables or settings.
    session_secret = settings.get('session.secret', os.environ.get('SESSION_SECRET'))
    
    if not session_secret:
        # CRITICAL: In production, this should fail or use a secure vault.
        # Using a random secret for development if none is provided.
        session_secret = os.urandom(32).hex()
        print("WARNING: Using ephemeral session secret. Configure 'session.secret' for production.")
        
    session_factory = SignedCookieSessionFactory(
        session_secret.encode('utf-8'),
        httponly=True, # Prevent client-side JS access to the cookie
        secure=True,   # Only send over HTTPS (requires proper deployment setup)
        timeout=3600   # Session timeout in seconds
    )

    with Configurator(settings=settings) as config:
        
        # 1. Security Configuration
        config.set_session_factory(session_factory)
        
        # 2. Database Initialization
        # Initialize DB engine and session factory based on settings
        initialize_db(settings)
        
        # 3. Include Routes and Views
        config.include('.routes')
        config.scan('.views')
        config.scan('.models')
        
    return config.make_wsgi_app()