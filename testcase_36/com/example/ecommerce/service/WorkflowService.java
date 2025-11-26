package com.example.ecommerce.service;

import com.example.ecommerce.enums.OrderState;
import com.example.ecommerce.model.Order;
import com.example.ecommerce.repository.OrderRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.util.NoSuchElementException;

@Service
public class WorkflowService {

    private final OrderRepository orderRepository;

    @Autowired
    public WorkflowService(OrderRepository orderRepository) {
        this.orderRepository = orderRepository;
    }

    /**
     * Step 1: Initialize a new order in the CART state.
     */
    @Transactional
    public Order initiateCart(BigDecimal amount, String email) {
        // Input validation is primarily handled by the Controller, but defensive checks remain good practice.
        if (amount == null || amount.compareTo(BigDecimal.ZERO) <= 0) {
            throw new IllegalArgumentException("Order amount must be positive.");
        }

        Order newOrder = new Order(amount, email);
        return orderRepository.save(newOrder);
    }

    /**
     * Step 2: Move order from CART to PAID_PENDING.
     * Simulates successful payment processing.
     */
    @Transactional
    public Order processPayment(Long orderId) {
        // Use parameterized query via JPA repository to prevent SQL Injection
        Order order = orderRepository.findById(orderId)
                .orElseThrow(() -> new NoSuchElementException("Order not found for ID: " + orderId));

        // SECURITY CHECK: Enforce sequential transition (CART -> PAID_PENDING)
        if (order.getCurrentStatus() != OrderState.CART) {
            throw new IllegalStateException(
                    String.format("Cannot process payment. Order %d is currently in state: %s",
                            orderId, order.getCurrentStatus()));
        }

        // Simulate external payment gateway interaction success
        order.setCurrentStatus(OrderState.PAID_PENDING);
        return orderRepository.save(order);
    }

    /**
     * Step 3: Finalize the order, moving it from PAID_PENDING to CONFIRMED.
     *
     * VULNERABILITY PREVENTION (CWE-841): State Transition Bypass
     * This method implements the mandatory check to ensure payment occurred.
     */
    @Transactional
    public Order finalizeOrder(Long orderId) {
        // Use parameterized query via JPA repository to prevent SQL Injection
        Order order = orderRepository.findById(orderId)
                .orElseThrow(() -> new NoSuchElementException("Order not found for ID: " + orderId));

        // --- CRITICAL SECURE WORKFLOW ENFORCEMENT (CWE-841 Mitigation) ---
        // Check if the order is in the required prerequisite state (PAID_PENDING).
        // This prevents an attacker from skipping the payment step and confirming the order directly.
        // VULNERABILITY INJECTED: The state check is removed, allowing confirmation regardless of payment status.
        
        // --- END SECURE CHECK ---

        // Critical business logic: Confirmation
        order.setCurrentStatus(OrderState.CONFIRMED);

        return orderRepository.save(order);
    }
}