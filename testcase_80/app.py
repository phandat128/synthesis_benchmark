import os
from flask import Flask, render_template, redirect, url_for, flash
from config import Config
from routes.checkout import checkout

# Ensure the templates directory is available for the 404 handler

def create_app():
    app = Flask(__name__)
    app.config.from_object(Config)

    # Security check: Warn if using the default insecure key
    if 'default_secret_key' in app.config['SECRET_KEY']:
        print("\n!!! SECURITY WARNING: Using default insecure secret key. Change SECRET_KEY in production. !!!\n")
        
    # Register blueprints
    app.register_blueprint(checkout)

    @app.errorhandler(404)
    def page_not_found(e):
        # Proper Error Handling: Do not leak sensitive information on error
        return render_template('404.html'), 404

    @app.route('/')
    def index():
        # Redirect to the start of the checkout process
        return redirect(url_for('checkout.cart_view'))

    return app

if __name__ == '__main__':
    app = create_app()
    # In production, use a WSGI server and enforce HTTPS/SSL termination.
    app.run(debug=app.config['DEBUG'])