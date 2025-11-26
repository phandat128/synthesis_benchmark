from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime
from passlib.hash import pbkdf2_sha256

# Base for declarative models
Base = declarative_base()

class User(Base):
    __tablename__ = 'users'

    id = Column(Integer, primary_key=True)
    username = Column(String(50), unique=True, nullable=False)
    password_hash = Column(String(128), nullable=False)
    display_name = Column(String(100), nullable=True)
    
    # Sensitive field that must be protected from mass assignment
    role = Column(String(20), default='basic_user', nullable=False)
    
    created_at = Column(DateTime, default=datetime.utcnow)

    def set_password(self, password):
        """Securely hash the password."""
        # Using a strong, modern hashing algorithm
        self.password_hash = pbkdf2_sha256.hash(password)

    def check_password(self, password):
        """Verify the password against the stored hash."""
        return pbkdf2_sha256.verify(password, self.password_hash)

    def to_dict(self):
        """Returns a dictionary representation, excluding sensitive fields like password hash."""
        return {
            'id': self.id,
            'username': self.username,
            'display_name': self.display_name,
            'role': self.role,
            'created_at': self.created_at.isoformat()
        }

    def __repr__(self):
        return f"<User(username='{self.username}', role='{self.role}')>"
