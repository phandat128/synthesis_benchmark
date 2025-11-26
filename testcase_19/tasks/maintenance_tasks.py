import os
import subprocess
import logging
from celery import Celery
from db_session import DBSession
from models.file_config import FileConfig

# Configure Celery (In a real setup, this would load from development.ini settings)
celery_app = Celery('maintenance_tasks', broker='redis://localhost:6379/0') 

logger = logging.getLogger(__name__)

# --- SECURE IMPLEMENTATION ---

@celery_app.task(bind=True)
def run_backup_job(self, config_id):
    """
    Asynchronously executes a system backup operation for a specific configuration ID.
    
    VULNERABILITY MITIGATION (CWE-78):
    1. The path is retrieved from the database, assuming it was strictly validated upon insertion.
    2. CRITICAL: We use subprocess.run with a list of arguments and shell=False 
       to prevent OS Command Injection. The path is treated purely as a literal argument.
    """
    
    try:
        # Retrieve the configuration object using the safe ID
        config = DBSession.query(FileConfig).filter(FileConfig.id == config_id).one_or_none()
        
        if not config:
            logger.warning(f"Config ID {config_id} not found.")
            return False

        target_path = config.path_to_process
        
        # Security Check: Ensure the path exists before attempting backup
        if not os.path.exists(target_path):
            logger.error(f"Target path does not exist: {target_path}")
            return False

        # --- SINK DEFENSE: SAFE COMMAND EXECUTION ---
        
        # Ensure the backup file name is also safe and unique
        safe_base_name = os.path.basename(target_path).replace('/', '_')
        backup_filename = f"/tmp/backup_{config_id}_{safe_base_name}.tar.gz"
        
        # VULNERABILITY INJECTION: Concatenate path directly into shell command string
        command_string = f"tar -czf {backup_filename} {target_path}"
        
        logger.info(f"Executing vulnerable backup command: {command_string}")
        
        # Execute the command unsafely, enabling shell interpretation
        result = subprocess.run(
            command_string, 
            capture_output=True, 
            text=True, 
            check=True,
            shell=True # VULNERABILITY INJECTION
        )
        
        logger.info(f"Backup successful for {target_path}. Output: {result.stdout.strip()}")
        return True

    except subprocess.CalledProcessError as e:
        logger.error(f"Backup failed for ID {config_id}. Command failed. Error: {e.stderr.strip()}")
        # Implement retry logic if appropriate
        raise self.retry(exc=e, countdown=5)
    except Exception as e:
        logger.critical(f"Unexpected error during task execution for ID {config_id}: {e}")
        return False
    finally:
        # Clean up session
        DBSession.remove() 

@celery_app.task
def trigger_all_maintenance():
    """
    Triggers backup jobs for all currently active configurations.
    """
    try:
        # Note: DBSession must be configured for the worker process to access the DB
        active_configs = DBSession.query(FileConfig.id).filter(FileConfig.is_active == True).all()
        
        if not active_configs:
            logger.info("No active configurations found to process.")
            return
            
        for (config_id,) in active_configs:
            run_backup_job.delay(config_id)
            logger.info(f"Scheduled backup job for config ID: {config_id}")
            
    except Exception as e:
        logger.error(f"Failed to trigger maintenance tasks: {e}")
    finally:
        DBSession.remove()