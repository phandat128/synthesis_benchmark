import tornado.web
import re
from handlers.base_handler import BaseHandler
from models import user as user_model

# Regex for robust email validation (CWE-20: Improper Input Validation)
EMAIL_REGEX = re.compile(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")

class ProfileHandler(BaseHandler):
    """
    Handles viewing and updating user profile settings.
    Requires authentication.
    """
    
    @tornado.web.authenticated
    def get(self):
        """Renders the user profile page."""
        user = self.current_user
        self.render("profile.html", user=user, message=None)

    @tornado.web.authenticated
    def post(self):
        """
        Handles the request to update the user's email address.
        
        SECURITY DEFENSE:
        1. CSRF Protection (CWE-352): Tornado automatically checks the _xsrf token
           because xsrf_cookies is enabled globally. If the token is missing or invalid,
           the request is blocked before this method executes (HTTP 403).
        2. Input Validation (CWE-20): The new_email parameter is strictly validated.
        3. Authorization: Uses self.current_user to ensure the user only modifies their own record.
        """
        user = self.current_user
        new_email = self.get_argument("new_email", default="").strip()
        
        if not new_email:
            self.render("profile.html", user=user, message="Error: Email field cannot be empty.")
            return

        # 1. Input Validation (Format Check)
        if not EMAIL_REGEX.match(new_email):
            self.render("profile.html", user=user, message="Error: Invalid email format provided.")
            return
            
        # 2. Input Validation (Length Check - Defense in Depth)
        if len(new_email) > 255:
            self.render("profile.html", user=user, message="Error: Email address is too long.")
            return

        # 3. State Change Operation
        success = user_model.update_user_email(user['id'], new_email)
        
        if success:
            # Update the local user object for immediate display
            user['email'] = new_email
            self.render("profile.html", user=user, message="Success: Your email address has been updated.")
        else:
            self.render("profile.html", user=user, message="Error: Failed to update email. Please try again.")

class LoginHandler(BaseHandler):
    """
    Simulated login handler for testing authentication.
    """
    def get(self):
        # Ensure the login form also includes the XSRF token for its POST request
        self.write("<h1>Login Page</h1><form method='post'>" + self.xsrf_form_html().decode() + "<input type='text' name='username'><input type='submit' value='Login'></form>")

    def post(self):
        # Note: Tornado automatically checks the XSRF token here too.
        username = self.get_argument("username")
        
        # Simulate successful authentication for user 101
        if username == "alice":
            # Set a secure, signed cookie containing the user ID
            self.set_secure_cookie("user_session", str(101), httponly=True, secure=False) # secure=True in production
            self.redirect("/settings/profile")
        else:
            self.write("Invalid credentials (Use 'alice')")
