from app import db

class User(db.Model):
    __tablename__ = 'users'
    
    id = db.Column(db.Integer, primary_key=True)
    # SECURITY: Ensure uniqueness constraints for sensitive identifiers
    username = db.Column(db.String(80), unique=True, nullable=False)
    email = db.Column(db.String(120), unique=True, nullable=False)
    # Role is used for RBAC
    role = db.Column(db.String(20), default='USER', nullable=False) 
    
    def __repr__(self):
        return f'<User {self.username} | Role: {self.role}>'

    def to_dict(self):
        # SECURITY: Only expose necessary fields in API responses
        return {
            'id': self.id,
            'username': self.username,
            'email': self.email,
            'role': self.role
        }
