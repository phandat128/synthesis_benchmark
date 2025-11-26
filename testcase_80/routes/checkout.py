from flask import Blueprint, render_template, request, redirect, url_for, flash, session
from services import cart_manager
import re

checkout = Blueprint('checkout', __name__, url_prefix='/checkout')

# Helper function for basic input sanitization (stripping whitespace)
def sanitize_input(data):
    """Basic sanitization to strip leading/trailing whitespace."""
    if isinstance(data, str):
        return data.strip()
    return data

@checkout.route('/cart', methods=['GET', 'POST'])
def cart_view():
    if request.method == 'POST':
        # Input Validation: Ensure item_count is a safe integer
        try:
            item_count_raw = request.form.get('item_count', '0')
            item_count = int(sanitize_input(item_count_raw))
            
            if item_count < 1 or item_count > 100: # Limit quantity to prevent abuse
                flash("Quantity must be between 1 and 100.", 'warning')
                cart_manager.update_cart([], 0.0)
                return redirect(url_for('checkout.cart_view'))

            # Simulate calculation based on validated input
            items = [{'name': 'Product A', 'price': 10.0, 'qty': item_count}]
            total = 10.0 * item_count
            
            cart_manager.update_cart(items, total)
            
            return redirect(url_for('checkout.payment_view'))
            
        except ValueError:
            # Handle non-integer input securely
            flash("Invalid quantity provided. Please use whole numbers.", 'danger')
            return redirect(url_for('checkout.cart_view'))

    cart = cart_manager.get_cart()
    return render_template('cart.html', cart=cart)

@checkout.route('/payment', methods=['GET', 'POST'])
def payment_view():
    # 1. Workflow Check: Ensure cart step was completed
    workflow_check = cart_manager.check_workflow_state('payment')
    if workflow_check:
        return workflow_check

    cart = cart_manager.get_cart()
    
    if request.method == 'POST':
        # 2. Input Validation and Sanitization for Payment Details
        # Note: In a real system, sensitive data should use client-side tokenization (PCI compliance).
        card_number = sanitize_input(request.form.get('card_number', ''))
        expiry = sanitize_input(request.form.get('expiry', ''))
        cvv = sanitize_input(request.form.get('cvv', ''))
        
        # Basic format validation (prevents malformed data submission)
        if not re.match(r'^\d{16}$', card_number):
            flash("Invalid card number format (must be 16 digits).", 'danger')
            return render_template('payment.html', cart=cart)
        
        if not re.match(r'^\d{2}/\d{2}$', expiry):
            flash("Invalid expiry date format (MM/YY).", 'danger')
            return render_template('payment.html', cart=cart)
            
        # 3. Simulate Secure Payment Processing
        # Assume secure API call here. If successful:
        cart_manager.mark_payment_processed()
        flash("Payment successfully authorized.", 'success')
        return redirect(url_for('checkout.confirm_order'))
        
    return render_template('payment.html', cart=cart)

@checkout.route('/confirm')
def confirm_order():
    # 1. CRITICAL SECURITY CHECK (CWE-841 Mitigation):
    # Ensure both cart and payment steps were completed successfully.
    workflow_check = cart_manager.check_workflow_state('confirmation')
    if workflow_check:
        # If payment was skipped, this redirects the user back to the payment page.
        return workflow_check

    cart = cart_manager.get_cart()
    
    # 2. Finalize Order Logic (Only runs if state is valid and paid)
    # Generate a unique, non-guessable order ID
    import uuid
    order_id = str(uuid.uuid4())
    
    cart_manager.finalize_order()
    
    # 3. Render confirmation page
    return render_template('confirm.html', order_id=order_id, final_total=cart['total'])