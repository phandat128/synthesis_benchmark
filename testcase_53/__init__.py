from pyramid.config import Configurator
from sqlalchemy import engine_from_config
from .models import initialize_sql, DBSession
from pyramid.session import SignedCookieSessionFactory
import os

def main(global_config, **settings):
    """
    This function returns a Pyramid WSGI application.
    """
    
    # --- 1. Database Configuration ---
    # Using SQLite for a self-contained example, but configuration is ready for PostgreSQL/psycopg2.
    if 'sqlalchemy.url' not in settings:
        # In a production environment, this should be a secure external DB connection string.
        settings['sqlalchemy.url'] = 'sqlite:///./comments.sqlite'
        
    engine = engine_from_config(settings, 'sqlalchemy.')
    
    # Initialize DB structure
    initialize_sql(engine)
    
    # Bind the session to the engine
    DBSession.configure(bind=engine)
    
    # --- 2. Session Configuration (Security Best Practice) ---
    # Use a strong, randomly generated secret key for signed cookies
    # In production, this MUST be read from a secure environment variable.
    session_secret = os.environ.get('SESSION_SECRET', 'a_very_long_and_secure_default_key_for_testing_1234567890_to_be_replaced')
    session_factory = SignedCookieSessionFactory(
        session_secret,
        # Secure flag ensures cookies are only sent over HTTPS (critical in production)
        secure=True, 
        # HttpOnly flag prevents client-side script access to the cookie (XSS mitigation)
        httponly=True,
        # Set max age for session validity
        max_age=86400 # 24 hours
    )

    # --- 3. Pyramid Configuration ---
    with Configurator(settings=settings) as config:
        
        # Set the session factory
        config.set_session_factory(session_factory)
        
        # Add transaction manager integration (for SQLAlchemy)
        config.include('pyramid_tm')
        
        # Setup Jinja2 templating
        config.include('pyramid_jinja2')
        config.add_jinja2_renderer('.jinja2')
        
        # Scan for view callables
        config.scan('.views')
        
        # --- 4. Route Configuration ---
        config.add_route('home', '/')
        config.add_route('submit', '/submit_comment')
        
    return config.make_wsgi_app()