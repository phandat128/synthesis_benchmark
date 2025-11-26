from app import db
from sqlalchemy import Integer, String, Text, Column, ForeignKey
from sqlalchemy.orm import relationship

class Document(db.Model):
    __tablename__ = 'documents'
    
    id = Column(Integer, primary_key=True)
    title = Column(String(255), nullable=False)
    content = Column(Text, nullable=False)
    
    # Link to the owner for auditing/ownership checks
    owner_id = Column(Integer, ForeignKey('users.id'), nullable=False)
    
    owner = relationship("User")

    def __repr__(self):
        return f'<Document {self.title}>'
        
    def to_dict(self):
        # SECURITY: Ensure no sensitive metadata (like internal file paths) is exposed.
        return {
            'id': self.id,
            'title': self.title,
            'content': self.content,
            'owner_id': self.owner_id
        }