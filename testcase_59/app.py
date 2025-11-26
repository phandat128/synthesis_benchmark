from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_session import Session
from config import Config
import os

# Initialize extensions globally
db = SQLAlchemy()
sess = Session()

def create_app():
    app = Flask(__name__)
    app.config.from_object(Config)

    # Initialize extensions
    db.init_app(app)
    sess.init_app(app)

    # Import and register blueprints
    from routes.document_routes import document_bp
    app.register_blueprint(document_bp, url_prefix='/api/v1')

    # Setup database (in a real app, this would be handled by migrations)
    with app.app_context():
        # Import models so SQLAlchemy knows about them
        import models.user
        import models.document
        db.create_all()
        
        # --- SECURE INITIALIZATION: Add mock data for testing authorization ---
        if not models.user.User.query.first():
            from services.auth_service import hash_password
            
            # User 1 (Alice): Authorized (Has both required groups)
            user_a = models.user.User(id=1, username='alice', password_hash=hash_password('securepass'), groups=['CONFIDENTIAL_ACCESS', 'PROJECT_ALPHA_LEAD', 'GENERAL_USER'])
            # User 2 (Bob): Unauthorized (Missing one group - tests the defense against CWE-863)
            user_b = models.user.User(id=2, username='bob', password_hash=hash_password('securepass'), groups=['CONFIDENTIAL_ACCESS', 'GENERAL_USER'])
            # User 3 (Charlie): Unauthorized (Missing both groups)
            user_c = models.user.User(id=3, username='charlie', password_hash=hash_password('securepass'), groups=['GENERAL_USER'])
            
            db.session.add_all([user_a, user_b, user_c])
            
            # Add a document (ID 1)
            doc = models.document.Document(id=1, title='Top Secret Project Alpha Report', content='This content is highly restricted and requires dual authorization.', owner_id=1)
            db.session.add(doc)
            
            db.session.commit()
            print("Database initialized with mock users and documents.")


    return app

if __name__ == '__main__':
    # SECURITY WARNING: Ensure SECRET_KEY is set via environment variables in production
    if Config.SECRET_KEY == 'default_insecure_key_CHANGE_ME':
        print("WARNING: Using default insecure secret key! Change config.py.")
        
    app = create_app()
    app.run(debug=True)