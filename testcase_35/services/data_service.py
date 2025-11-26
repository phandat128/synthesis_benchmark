from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine
from models.customer import Customer

# Global session factory placeholder
DBSession = None

def initialize_db(settings):
    """Initializes the database engine and session factory."""
    # Use a secure connection string from settings
    db_url = settings.get('sqlalchemy.url', 'sqlite:///./default.db')
    engine = create_engine(db_url)
    
    # Ensure tables exist (for demonstration)
    # Customer.Base.metadata.create_all(engine)
    
    global DBSession
    DBSession = sessionmaker(bind=engine)
    return engine

def fetch_customer_data(limit: int) -> list:
    """
    Fetches customer records up to the specified limit.
    
    SECURITY NOTE: The 'limit' parameter MUST be strictly validated 
    by the calling view function (views/report_views.py) to prevent 
    resource exhaustion (CWE-400) before reaching this sink.
    
    :param limit: The maximum number of records to fetch (must be safe).
    :return: A list of customer dictionaries.
    """
    if DBSession is None:
        raise RuntimeError("Database session not initialized.")

    session = DBSession()
    try:
        # Use SQLAlchemy's ORM methods (.limit()), which are inherently safe 
        # against SQL injection. The security focus here is on the value of 'limit'.
        customers = session.query(Customer).limit(limit).all()
        
        # Convert ORM objects to dictionaries for transport/processing
        return [c.to_dict() for c in customers]
    except Exception as e:
        # Log the error internally (not shown)
        print(f"Database error during fetch: {e}")
        # Raise a generic, non-sensitive error for the caller
        raise ConnectionError("Could not retrieve customer data.")
    finally:
        session.close()