from fastapi import APIRouter, Depends, status, HTTPException
from models import Cart, PaymentDetails, CheckoutStatus
from services.checkout_service import checkout_service, CheckoutService
import uuid

router = APIRouter(
    prefix="/api/checkout",
    tags=["checkout"],
)

# Helper function for path validation
def validate_order_id(order_id: str):
    """Ensures the order_id is a valid UUID format to prevent injection/malformed input."""
    try:
        # Strict validation of the expected format
        uuid.UUID(order_id)
    except ValueError:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Invalid order ID format."
        )
    return order_id

@router.post("/", response_model=CheckoutStatus, status_code=status.HTTP_201_CREATED)
def initiate_checkout(service: CheckoutService = Depends(lambda: checkout_service)):
    """Starts a new checkout session and returns the order ID."""
    return service.start_checkout()

@router.post("/{order_id}/cart", response_model=CheckoutStatus)
def update_order_cart(
    order_id: str = Depends(validate_order_id),
    cart: Cart,
    service: CheckoutService = Depends(lambda: checkout_service)
):
    """Step 1: Update the shopping cart details. Input is validated by Pydantic (Cart schema)."""
    return service.update_cart(order_id, cart)

@router.post("/{order_id}/payment", response_model=CheckoutStatus)
def process_order_payment(
    order_id: str = Depends(validate_order_id),
    payment: PaymentDetails,
    service: CheckoutService = Depends(lambda: checkout_service)
):
    """Step 2: Process payment details. Input is validated by Pydantic (PaymentDetails schema)."""
    return service.process_payment(order_id, payment)

@router.post("/{order_id}/confirm", response_model=CheckoutStatus)
def confirm_order(
    order_id: str = Depends(validate_order_id),
    service: CheckoutService = Depends(lambda: checkout_service)
):
    """
    Step 3: Final confirmation.
    The service layer strictly enforces the preceding PAYMENT_PROCESSED state (CWE-841 mitigation).
    """
    return service.finalize_order(order_id)

@router.get("/{order_id}", response_model=CheckoutStatus)
def get_order_status(
    order_id: str = Depends(validate_order_id),
    service: CheckoutService = Depends(lambda: checkout_service)
):
    """Retrieves the current status of an order."""
    status = service._get_and_validate_order(order_id)
    # Ensure sensitive payment data is explicitly removed before returning to a general endpoint
    status.payment_data = None 
    return status