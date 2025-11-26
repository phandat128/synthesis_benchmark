import tornado.ioloop
import tornado.web
import os

from config import Config
from handlers.user_handler import ProfileHandler, LoginHandler
from handlers.base_handler import BaseHandler

# Define a simple logout handler for completeness
class LogoutHandler(BaseHandler):
    def get(self):
        self.clear_cookie("user_session")
        self.redirect("/login")

def make_app():
    """
    Initializes the Tornado application with secure settings and routing.
    """
    settings = {
        "template_path": os.path.join(os.path.dirname(__file__), "templates"),
        "static_path": os.path.join(os.path.dirname(__file__), "static"),
        "cookie_secret": Config.COOKIE_SECRET,
        "login_url": "/login",
        
        # CRITICAL SECURITY SETTING: Enables XSRF protection globally.
        # Tornado automatically checks for the _xsrf token on POST/PUT/DELETE requests
        # if this is True, preventing Cross-Site Request Forgery (CWE-352).
        "xsrf_cookies": Config.XSRF_COOKIES, 
        
        "debug": Config.DEBUG,
    }

    return tornado.web.Application([
        (r"/login", LoginHandler),
        (r"/logout", LogoutHandler),
        (r"/settings/profile", ProfileHandler), # GET request to view profile
        # Route for the state-changing operation (POST protected by CSRF)
        (r"/settings/email", ProfileHandler), 
    ], **settings)

if __name__ == "__main__":
    app = make_app()
    print(f"Starting server on http://127.0.0.1:{Config.PORT}")
    print(f"CSRF Protection Status: {'Enabled' if Config.XSRF_COOKIES else 'Disabled'}")
    app.listen(Config.PORT)
    tornado.ioloop.IOLoop.current().start()
