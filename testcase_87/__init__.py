import os
from pyramid.config import Configurator
from sqlalchemy import engine_from_config
from sqlalchemy.orm import sessionmaker
from pyramid.authentication import AuthTktAuthenticationPolicy
from pyramid.authorization import ACLAuthorizationPolicy
from .security.auth_policy import JWTAuthPolicy, groupfinder

# Database setup
Session = sessionmaker()

def db_session(request):
    """Provide a database session object to the request."""
    maker = request.registry.dbmaker
    session = maker()
    
    # Use pyramid_tm for transaction management
    def cleanup(request):
        if request.exception is not None:
            session.rollback()
        else:
            session.commit()
        session.close()

    request.add_finished_callback(cleanup)
    return session


def main(global_config, **settings):
    """ This function returns a Pyramid WSGI application. """
    
    # 1. Initialize Configurator
    with Configurator(settings=settings) as config:
        
        # 2. Database Configuration
        engine = engine_from_config(settings, 'sqlalchemy.')
        config.registry.dbmaker = sessionmaker(bind=engine)
        
        # Import models and create tables (for development/setup)
        from .models.user import Base
        Base.metadata.create_all(engine)
        
        # Add the database session to the request object
        config.add_request_method(db_session, 'db', reify=True)

        # 3. Security Configuration (using JWT)
        # The JWTAuthPolicy handles token verification and extraction
        config.set_authentication_policy(
            JWTAuthPolicy(
                secret=settings['jwt.secret'],
                algorithm=settings['jwt.algorithm'],
                callback=groupfinder
            )
        )
        config.set_authorization_policy(ACLAuthorizationPolicy())
        
        # 4. Include necessary packages
        config.include('pyramid_tm') # Transaction management
        config.include('.views')
        
        # 5. Routing
        config.add_route('profile_update', '/api/v1/profile', request_method='PUT')
        
        config.scan()

    return config.make_wsgi_app()
