from typing import Dict, Optional
from models import CheckoutStatus, CheckoutState, Cart, PaymentDetails
import uuid
from datetime import datetime

# In-memory store simulating a database table or persistent session store
# Key: order_id (str), Value: CheckoutStatus object
_STATE_STORE: Dict[str, CheckoutStatus] = {}

class StateStore:
    """
    Manages the persistence and retrieval of checkout state.
    Uses an in-memory dictionary for simplicity, simulating secure data storage.
    """

    @staticmethod
    def create_new_checkout() -> CheckoutStatus:
        """Creates a new order ID and initializes the state."""
        order_id = str(uuid.uuid4())
        now = datetime.utcnow().isoformat()
        
        new_status = CheckoutStatus(
            order_id=order_id,
            current_state=CheckoutState.INITIATED,
            last_updated=now
        )
        _STATE_STORE[order_id] = new_status
        return new_status

    @staticmethod
    def get_status(order_id: str) -> Optional[CheckoutStatus]:
        """Retrieves the current status for a given order ID."""
        # Safe dictionary access
        return _STATE_STORE.get(order_id)

    @staticmethod
    def update_status(order_id: str, new_state: CheckoutState, data_payload: Optional[dict] = None) -> CheckoutStatus:
        """Updates the state and optionally merges new data."""
        if order_id not in _STATE_STORE:
            # This should ideally be caught by the service layer, but defensive coding here.
            raise ValueError(f"Order ID {order_id} not found.")
        
        status = _STATE_STORE[order_id]
        
        # Update state and timestamp
        status.current_state = new_state
        status.last_updated = datetime.utcnow().isoformat()
        
        # Merge data payload if provided
        if data_payload:
            if 'cart_data' in data_payload and isinstance(data_payload['cart_data'], Cart):
                status.cart_data = data_payload['cart_data']
            if 'payment_data' in data_payload and isinstance(data_payload['payment_data'], PaymentDetails):
                status.payment_data = data_payload['payment_data']
                
        # Persist change (in-memory update is sufficient here)
        return status

    @staticmethod
    def clear_store():
        """Utility function for testing/cleanup."""
        _STATE_STORE.clear()