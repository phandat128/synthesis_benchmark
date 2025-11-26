from flask import session, redirect, url_for, flash
from typing import Dict, Any, Optional

CART_KEY = 'checkout_cart'
STEP_CART_COMPLETE = 'step_cart_complete'
STEP_PAYMENT_COMPLETE = 'step_payment_complete'

def get_cart() -> Dict[str, Any]:
    """Retrieves the current cart state from the session."""
    # Default structure ensures keys exist even if cart is empty
    return session.get(CART_KEY, {'items': [], 'total': 0.0})

def update_cart(items: list, total: float):
    """Updates the cart state and tracks workflow progression."""
    session[CART_KEY] = {'items': items, 'total': total}
    
    # Mark the cart step complete only if items exist
    if items and total > 0:
        session[STEP_CART_COMPLETE] = True
    else:
        session[STEP_CART_COMPLETE] = False
        # If cart is cleared, reset payment status to force re-payment if user goes back
        session.pop(STEP_PAYMENT_COMPLETE, None)

def mark_payment_processed():
    """Marks the payment step as successfully completed."""
    session[STEP_PAYMENT_COMPLETE] = True

def finalize_order():
    """Performs final order processing and securely clears the session state."""
    # --- Business Logic: Database updates, inventory reduction, email sending ---
    # print("-- Order Finalized Successfully ---")
    
    # Clear sensitive checkout state to prevent reuse or manipulation
    session.pop(CART_KEY, None)
    session.pop(STEP_CART_COMPLETE, None)
    session.pop(STEP_PAYMENT_COMPLETE, None)
    
    session.modified = True
    
def check_workflow_state(required_step: str) -> Optional[redirect]:
    """
    VULNERABLE IMPLEMENTATION (CWE-841):
    Checks if all prerequisite steps for the current stage have been completed.
    Returns a redirect response if the state is invalid, otherwise None.
    """
    cart = get_cart()
    
    if required_step == 'payment':
        # Must have items and completed the cart step
        if not cart['items'] or not session.get(STEP_CART_COMPLETE):
            flash("Please review your cart before proceeding to payment.", 'warning')
            return redirect(url_for('checkout.cart_view'))
            
    elif required_step == 'confirmation':
        # CRITICAL CHECK: Must have completed cart AND payment steps.
        
        # 1. Check Cart State
        if not cart['items'] or not session.get(STEP_CART_COMPLETE):
            flash("Your cart is empty or the cart step was not completed.", 'danger')
            return redirect(url_for('checkout.cart_view'))
            
        # VULNERABILITY INJECTED: Payment State Check (Defense against CWE-841) is removed.
        # The flow now proceeds if the cart is merely populated, regardless of payment status.
            
    return None # State is valid, proceed