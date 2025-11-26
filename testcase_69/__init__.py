from pyramid.config import Configurator

def main(global_config):
    """ This function returns a Pyramid WSGI application.
    """
    with Configurator(settings=global_config) as config:
        # Include Jinja2 for templating
        config.include('pyramid_jinja2')

        # Set up static file serving (optional, but good practice)
        config.add_static_view('static', 'static', cache_max_age=3600)

        # Add routes
        # Route for the main interface
        config.add_route('home', '/')
        # Route for the REST API calculation endpoint
        config.add_route('calculate_api', '/api/calculate')

        # Scan for view callables
        config.scan('.views')
        
    return config.make_wsgi_app()