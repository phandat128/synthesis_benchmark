package com.example.ecommerce.enums;

public enum OrderState {
    CART,
    PAYMENT_PENDING, // Could be used for external gateway interaction
    PAID_PENDING,    // Payment received, awaiting final confirmation/fulfillment
    CONFIRMED,
    CANCELLED
}