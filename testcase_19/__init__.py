from pyramid.config import Configurator
from pyramid.session import SignedCookieSessionFactory
from db_session import initialize_sql
import transaction

def main(global_config, **settings):
    """
    This function returns a Pyramid WSGI application.
    """
    
    # Secure Session Configuration
    # Use a strong, long secret key and secure cookie settings
    session_secret = settings.get('session.secret', 'a_very_long_and_secure_default_secret_key_for_testing_1234567890')
    session_factory = SignedCookieSessionFactory(
        session_secret,
        httponly=True,
        # Enforce HTTPS in production by setting session.secure=true in .ini
        secure=settings.get('session.secure', 'true').lower() == 'true',
        max_age=3600 # 1 hour session timeout
    )

    with Configurator(settings=settings) as config:
        
        # 1. Security: Setup Session
        config.set_session_factory(session_factory)
        
        # 2. Setup Database and Transaction Manager
        engine = initialize_sql(settings)
        config.include('pyramid_tm') # Integrates transaction management
        
        # 3. Setup Routes
        config.add_route('register_path', '/api/v1/config/register')
        config.add_route('trigger_maintenance', '/api/v1/maintenance/run')
        config.add_route('get_config', '/api/v1/config/{config_id}')

        # 4. Scan and Register Views
        config.scan('.views')
        
        # 5. Initialize Database Models (For initial setup/development)
        # In production, use Alembic migrations instead of create_all.
        from models.file_config import Base
        Base.metadata.create_all(engine)

    return config.make_wsgi_app()