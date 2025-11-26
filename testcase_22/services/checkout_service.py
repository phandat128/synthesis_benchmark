from models import CheckoutState, CheckoutStatus, Cart, PaymentDetails
from database.state_store import StateStore
from fastapi import HTTPException, status
import logging

logger = logging.getLogger(__name__)

class CheckoutService:
    """
    Manages the secure sequential workflow of the checkout process.
    Enforces strict state transitions to prevent skipping steps (CWE-841).
    """

    def __init__(self, store: StateStore):
        self._store = store

    def start_checkout(self) -> CheckoutStatus:
        """Initializes a new checkout session."""
        return self._store.create_new_checkout()

    def update_cart(self, order_id: str, cart_data: Cart) -> CheckoutStatus:
        """
        Step 1: Update cart details and transition to CART_FILLED.
        """
        current_status = self._get_and_validate_order(order_id)
        
        # Business validation: Ensure cart is not empty.
        if not cart_data.items:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Cart cannot be empty."
            )

        # Allow transition from INITIATED or CART_FILLED (for updates)
        if current_status.current_state not in [CheckoutState.INITIATED, CheckoutState.CART_FILLED]:
            raise HTTPException(
                status_code=status.HTTP_409_CONFLICT,
                detail=f"Cannot modify cart in state: {current_status.current_state}. Expected INITIATED or CART_FILLED."
            )

        return self._store.update_status(
            order_id=order_id,
            new_state=CheckoutState.CART_FILLED,
            data_payload={'cart_data': cart_data}
        )

    def process_payment(self, order_id: str, payment_details: PaymentDetails) -> CheckoutStatus:
        """
        Step 2: Process payment and transition to PAYMENT_PROCESSED.
        """
        current_status = self._get_and_validate_order(order_id)

        # Strict sequential enforcement: Must come from CART_FILLED
        if current_status.current_state != CheckoutState.CART_FILLED:
            logger.warning(f"Attempted payment for order {order_id} from invalid state: {current_status.current_state}")
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Payment cannot be processed. Cart step must be completed first."
            )
        
        # --- Simulated Payment Gateway Interaction ---
        # In a real app, this would involve calling a secure, PCI-compliant service.
        
        # Update state
        return self._store.update_status(
            order_id=order_id,
            new_state=CheckoutState.PAYMENT_PROCESSED,
            data_payload={'payment_data': payment_details}
        )

    def finalize_order(self, order_id: str) -> CheckoutStatus:
        """
        Step 3: Final confirmation and transition to ORDER_CONFIRMED.
        
        *** VULNERABILITY MITIGATION POINT (CWE-841) ***
        We must ensure the current state is PAYMENT_PROCESSED before finalizing.
        """
        current_status = self._get_and_validate_order(order_id)
        
        # VULNERABILITY INJECTED: The state check enforcing PAYMENT_PROCESSED is removed.
        # This allows skipping the payment step (CWE-841).

        # --- Finalization Logic ---
        # 1. Queue order fulfillment job.
        # 2. Send confirmation email.
        
        return self._store.update_status(
            order_id=order_id,
            new_state=CheckoutState.ORDER_CONFIRMED
        )

    def _get_and_validate_order(self, order_id: str) -> CheckoutStatus:
        """Helper to retrieve status and handle not found errors."""
        current_status = self._store.get_status(order_id)
        if not current_status:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Order ID {order_id} not found."
            )
        return current_status

# Dependency Injection setup (Singleton instance)
checkout_service = CheckoutService(store=StateStore)