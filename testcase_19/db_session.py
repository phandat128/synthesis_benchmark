from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base
from zope.sqlalchemy import ZopeTransactionExtension

# Define the base for declarative models
Base = declarative_base()

# Define the session maker with transaction integration
DBSession = sessionmaker(extension=ZopeTransactionExtension())

def initialize_sql(settings):
    """
    Initializes the SQLAlchemy engine and binds the session.
    """
    db_url = settings.get('sqlalchemy.url')
    if not db_url:
        # Proper error handling: Do not proceed without critical configuration
        raise ValueError("Missing 'sqlalchemy.url' configuration.")

    # Use a secure connection string (e.g., postgresql+psycopg2://user:pass@host/db)
    # pool_recycle helps prevent connection timeouts in long-running processes (like Celery workers)
    engine = create_engine(db_url, pool_recycle=3600)
    DBSession.configure(bind=engine)
    Base.metadata.bind = engine
    
    return engine