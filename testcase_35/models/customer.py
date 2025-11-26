from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime

Base = declarative_base()

class Customer(Base):
    """SQLAlchemy model for customer data."""
    __tablename__ = 'customers'

    id = Column(Integer, primary_key=True)
    first_name = Column(String(100), nullable=False)
    last_name = Column(String(100), nullable=False)
    email = Column(String(255), unique=True, nullable=False)
    # Store balance in cents to avoid floating point issues
    account_balance = Column(Integer, default=0)
    created_at = Column(DateTime, default=datetime.utcnow)

    def __repr__(self):
        return f"<Customer(id={self.id}, email='{self.email}')>"

    def to_dict(self):
        """Securely serialize data for report generation."""
        return {
            'id': self.id,
            'first_name': self.first_name,
            'last_name': self.last_name,
            'email': self.email,
            # Format balance securely for display
            'balance': f"${self.account_balance / 100:.2f}", 
            'created_at': self.created_at.isoformat()
        }