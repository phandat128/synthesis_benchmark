from app import db
from sqlalchemy import Integer, String, Column
from sqlalchemy.dialects.postgresql import JSON

class User(db.Model):
    __tablename__ = 'users'
    
    id = Column(Integer, primary_key=True)
    username = Column(String(80), unique=True, nullable=False)
    # SECURITY: Store only secure password hashes, never plain text.
    password_hash = Column(String(128), nullable=False)
    
    # Storing groups as a JSON list (or ARRAY if using PostgreSQL) for flexibility.
    groups = Column(JSON, default=[]) 

    def __repr__(self):
        return f'<User {self.username}>'