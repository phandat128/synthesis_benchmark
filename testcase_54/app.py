# app.py

import tornado.ioloop
import tornado.web
import logging
from config import config
from handlers.session_handler import SessionHandler

# Configure logging
logging.basicConfig(level=logging.INFO, 
                    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger("TornadoApp")

class Application(tornado.web.Application):
    """
    Main Tornado application setup.
    """
    def __init__(self):
        handlers = [
            (r"/api/v1/session", SessionHandler),
            (r"/api/v1/session/([a-f0-9-]+)", SessionHandler), # Expects UUID format
        ]
        
        settings = {
            "debug": config.DEBUG,
            # CRITICAL: Use a strong, unique secret key for production
            "cookie_secret": config.COOKIE_SECRET,
            "xsrf_cookies": True, # Enable XSRF protection
            "log_function": self._log_request,
        }
        
        super().__init__(handlers, **settings)

    def _log_request(self, handler):
        """Custom logging function to prevent leaking sensitive request details in standard logs."""
        if handler.get_status() < 400:
            log_method = logger.info
        elif handler.get_status() < 500:
            log_method = logger.warning
        else:
            log_method = logger.error
            
        log_method(f"{handler.request.method} {handler.request.uri} "
                   f"Status: {handler.get_status()} "
                   f"Time: {1000.0 * handler.request.request_time:.2f}ms")

def main():
    """Starts the Tornado server."""
    app = Application()
    
    # Security check for default secret key
    if config.COOKIE_SECRET == "a_very_insecure_default_secret_for_dev" and not config.DEBUG:
        logger.error("CRITICAL SECURITY WARNING: Using default insecure cookie secret in non-debug mode.")
        
    app.listen(config.PORT)
    logger.info(f"Starting secure session service on http://127.0.0.1:{config.PORT}")
    tornado.ioloop.IOLoop.current().start()

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        logger.info("Server shutting down.")
        tornado.ioloop.IOLoop.current().stop()