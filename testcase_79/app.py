import os
from flask import Flask, jsonify
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from flask_jwt_extended import JWTManager
from config import Config

# Initialize extensions outside of create_app
db = SQLAlchemy()
migrate = Migrate()
jwt = JWTManager()

def create_app(config_class=Config):
    app = Flask(__name__)
    app.config.from_object(config_class)

    # Initialize extensions
    db.init_app(app)
    migrate.init_app(app, db)
    jwt.init_app(app)

    # Import models so they are registered with SQLAlchemy
    from models.user import User

    # Register Blueprints
    from routes.user_routes import user_bp
    app.register_blueprint(user_bp)

    # JWT Error Handlers (Security: Do not leak internal details)
    @jwt.unauthorized_loader
    def unauthorized_callback(callback):
        return jsonify({"msg": "Missing or invalid token"}), 401

    @jwt.invalid_token_loader
    def invalid_token_callback(error):
        return jsonify({"msg": "Signature verification failed"}), 422
        
    @jwt.expired_token_loader
    def expired_token_callback(jwt_header, jwt_payload):
        return jsonify({"msg": "The token has expired"}), 401

    # Simple root route for health check
    @app.route('/')
    def index():
        return jsonify({"status": "running", "version": "1.0"})

    return app

if __name__ == '__main__':
    app = create_app()
    
    # Example setup for initial database creation and default user (for testing RBAC)
    with app.app_context():
        db.create_all() 
        
        from models.user import User
        if not User.query.filter_by(username='admin').first():
            admin_user = User(username='admin', email='admin@example.com', role='ADMIN')
            db.session.add(admin_user)
            db.session.commit()
            print("Created default 'admin' user.")
            
    app.run(debug=True)
