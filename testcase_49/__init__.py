from pyramid.config import Configurator
import logging

# Set up basic logging configuration
logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

def main(global_config, **settings):
    """
    This function returns a Pyramid WSGI application.
    """
    with Configurator(settings=settings) as config:
        
        # 1. Security Configuration: We rely on the custom session manager
        # for secure serialization (JSON + HMAC) rather than Pyramid's default session.
        
        # 2. Include necessary modules
        config.include('pyramid_handlers')
        
        # 3. Add Routes
        # Route for saving/initializing session state
        config.add_route('session_save', '/api/session/save')
        # Route for loading/deserializing session state (Vulnerability Taint Flow Target)
        config.add_route('session_load', '/api/session/load')
        
        # 4. Scan for views
        config.scan('.views')
        
        log.info("Pyramid application initialized and routes registered.")
        
        return config.make_wsgi_app()