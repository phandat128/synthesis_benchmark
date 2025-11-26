from pydantic import BaseModel, Field
from enum import Enum
from typing import List, Optional

class CheckoutState(str, Enum):
    # Initial state
    INITIATED = "INITIATED"
    # Step 1 complete
    CART_FILLED = "CART_FILLED"
    # Step 2 complete (The required preceding state for confirmation)
    PAYMENT_PROCESSED = "PAYMENT_PROCESSED"
    # Final state
    ORDER_CONFIRMED = "ORDER_CONFIRMED"
    # Error state
    FAILED = "FAILED"

class CartItem(BaseModel):
    product_id: str = Field(..., min_length=1, max_length=50)
    quantity: int = Field(..., gt=0)
    price: float = Field(..., gt=0.0)

class Cart(BaseModel):
    items: List[CartItem]
    # Total amount must be calculated client-side and validated server-side in a real app,
    # but here we validate structure.
    total_amount: float = Field(..., gt=0.0)

class PaymentDetails(BaseModel):
    # Using a tokenized approach for security; raw card data should never be handled here.
    payment_token: str = Field(..., min_length=10)
    billing_address: str = Field(..., min_length=5)

class CheckoutStatus(BaseModel):
    order_id: str
    current_state: CheckoutState
    last_updated: str
    cart_data: Optional[Cart] = None
    # Payment data is optional and should be stripped for most public views
    payment_data: Optional[PaymentDetails] = None