from sqlalchemy import Column, Integer, Text, DateTime, func
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import scoped_session, sessionmaker
from zope.sqlalchemy import ZopeTransactionExtension

# Define the base for declarative models
Base = declarative_base()

# Configure the session maker for transactional context
# ZopeTransactionExtension ensures transactions are managed automatically by Pyramid's transaction manager
DBSession = scoped_session(sessionmaker(extension=ZopeTransactionExtension()))

class Comment(Base):
    """
    SQLAlchemy model for a user comment.
    """
    __tablename__ = 'comments'
    id = Column(Integer, primary_key=True)
    # Store content as Text to allow for longer comments
    content = Column(Text, nullable=False)
    # Author name
    author = Column(Text, nullable=False, default='Anonymous')
    # Timestamp for ordering and display
    created_at = Column(DateTime, default=func.now())

    def __repr__(self):
        return f"<Comment(id={self.id}, author='{self.author}')>"

# Function to initialize the database structure
def initialize_sql(engine):
    """
    Creates the database tables if they do not exist.
    """
    Base.metadata.bind = engine
    Base.metadata.create_all(engine)