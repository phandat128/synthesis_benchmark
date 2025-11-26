import logging
from pyramid.config import Configurator
from sqlalchemy import engine_from_config
from pyramid_report_app.models.data_model import DBSession, Base

log = logging.getLogger(__name__)


def get_authenticated_user(request):
    """
    Placeholder for a real authentication system (e.g., JWT, session).
    In a real application, this would securely retrieve the user identity
    and ensure least privilege principles are followed.
    """
    # Simulate a successful authentication for user ID 1
    # This dictionary represents the secure context of the currently logged-in user.
    return {'user_id': 1, 'username': 'reporter_user'}


def db_session_factory(request):
    """
    Provides a database session tied to the request lifecycle.
    """
    return DBSession


def main(global_config, **settings):
    """ This function returns a Pyramid WSGI application. """
    
    # 1. Database Setup
    engine = engine_from_config(settings, 'sqlalchemy.')
    DBSession.configure(bind=engine)
    
    # Ensure tables exist (for initial setup/testing)
    # In production, this is usually handled by migrations (e.g., Alembic)
    Base.metadata.create_all(engine)
    
    # 2. Pyramid Configuration
    with Configurator(settings=settings) as config:
        
        # Include necessary components
        config.include('pyramid_jinja2')
        config.include('pyramid_tm') # Transaction Manager for SQLAlchemy
        
        # Add the request property for the authenticated user (for secure context)
        config.add_request_method(get_authenticated_user, 'authenticated_user', reify=True)
        
        # Add the request property for the database session
        config.add_request_method(db_session_factory, 'dbsession', reify=True)

        # 3. Routing
        config.add_route('generate_report', '/api/v1/generate_report')
        
        # 4. Scanning for views
        config.scan('.views')
        config.scan('.models')
        
    log.info("Pyramid application initialized securely.")
    return config.make_wsgi_app()