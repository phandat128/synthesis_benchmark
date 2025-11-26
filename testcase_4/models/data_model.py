import datetime
from sqlalchemy import (
    Column,
    Integer,
    Text,
    DateTime,
    Boolean
)
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import scoped_session, sessionmaker
from zope.sqlalchemy import ZopeTransactionExtension

# Setup database session management
DBSession = scoped_session(sessionmaker(extension=ZopeTransactionExtension()))
Base = declarative_base()

class DataModel(Base):
    """
    Represents a record in the database that will be included in the report.
    """
    __tablename__ = 'report_data'

    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, nullable=False, index=True)
    report_title = Column(Text, nullable=False)
    content_hash = Column(Text, nullable=False) # Internal/sensitive data
    is_processed = Column(Boolean, default=False)
    created_at = Column(DateTime, default=datetime.datetime.utcnow)

    def __repr__(self):
        return f"<DataModel(id={self.id}, title='{self.report_title[:20]}...')>"

    def to_dict(self):
        """Securely serialize data for report generation, omitting sensitive fields."""
        return {
            'id': self.id,
            'user_id': self.user_id,
            'title': self.report_title,
            'created_at': self.created_at.isoformat()
            # Note: content_hash is sensitive/internal and is omitted from the report data structure
        }