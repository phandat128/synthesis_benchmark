import datetime
from sqlalchemy import (
    Column,
    Integer,
    String,
    DateTime,
    Boolean,
)
from sqlalchemy.exc import IntegrityError
from db_session import Base, DBSession

class FileConfig(Base):
    """
    SQLAlchemy model for storing file configuration paths requested by users.
    """
    __tablename__ = 'file_configurations'

    id = Column(Integer, primary_key=True)
    # path_to_process stores the validated, sanitized file path.
    path_to_process = Column(String(512), nullable=False, unique=True)
    user_id = Column(String(64), nullable=False) # Placeholder for user identification
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime, default=datetime.datetime.utcnow)

    @classmethod
    def create_new_config(cls, user_id, validated_path):
        """
        Creates and persists a new FileConfig entry.
        
        :param user_id: Identifier of the user requesting the configuration.
        :param validated_path: The path, already checked for safety and validity.
        :return: The newly created FileConfig object or None if creation failed.
        """
        try:
            new_config = cls(
                user_id=user_id,
                path_to_process=validated_path
            )
            DBSession.add(new_config)
            DBSession.flush() # Ensure the transaction sees the new object
            return new_config
        except IntegrityError:
            # Handle case where path_to_process (unique constraint) already exists
            DBSession.rollback()
            return None
        except Exception as e:
            # General error handling, log the error but don't leak details
            print(f"Database error during config creation: {e}")
            DBSession.rollback()
            return None

    @classmethod
    def get_active_paths(cls):
        """
        Retrieves all active file paths for maintenance tasks.
        """
        return DBSession.query(cls.path_to_process).filter(cls.is_active == True).all()