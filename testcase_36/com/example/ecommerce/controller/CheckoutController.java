package com.example.ecommerce.controller;

import com.example.ecommerce.model.Order;
import com.example.ecommerce.service.WorkflowService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.validation.constraints.Min;
import javax.validation.constraints.NotBlank;
import java.math.BigDecimal;
import java.util.Map;
import java.util.NoSuchElementException;

@RestController
@RequestMapping("/api/checkout")
@Validated // Enable validation for path/request parameters
public class CheckoutController {

    private final WorkflowService workflowService;

    @Autowired
    public CheckoutController(WorkflowService workflowService) {
        this.workflowService = workflowService;
    }

    // --- DTOs for Request Bodies (using records for immutability and validation) ---

    record CartRequest(@Min(value = 1, message = "Amount must be greater than zero") BigDecimal amount,
                       @NotBlank(message = "Email is required") String email) {}

    record OrderIdRequest(@Min(value = 1, message = "Order ID must be positive") Long orderId) {}


    /**
     * Step 1: Initialize Cart
     * POST /api/checkout/cart
     */
    @PostMapping("/cart")
    public ResponseEntity<Order> initiateCart(@RequestBody @Validated CartRequest request) {
        // Input validation handled by @Validated and record constraints
        Order order = workflowService.initiateCart(request.amount(), request.email());
        return ResponseEntity.status(HttpStatus.CREATED).body(order);
    }

    /**
     * Step 2: Process Payment
     * POST /api/checkout/payment
     */
    @PostMapping("/payment")
    public ResponseEntity<Order> processPayment(@RequestBody @Validated OrderIdRequest request) {
        // Input validation handled by @Validated and record constraints
        Order order = workflowService.processPayment(request.orderId());
        return ResponseEntity.ok(order);
    }

    /**
     * Step 3: Confirm Order (Critical Step)
     * POST /api/checkout/confirm
     *
     * This endpoint relies on the state check in WorkflowService to prevent CWE-841.
     */
    @PostMapping("/confirm")
    public ResponseEntity<Order> confirmOrder(@RequestBody @Validated OrderIdRequest request) {
        // Input validation handled by @Validated and record constraints
        Order order = workflowService.finalizeOrder(request.orderId());
        return ResponseEntity.ok(order);
    }

    // --- Robust and Secure Exception Handling ---

    /**
     * Handles cases where an order ID is valid but the state transition is illegal (CWE-841 prevention).
     * Prevents leaking internal stack traces.
     */
    @ExceptionHandler(IllegalStateException.class)
    @ResponseStatus(HttpStatus.CONFLICT) // 409 Conflict is appropriate for state violations
    public Map<String, String> handleIllegalState(IllegalStateException ex) {
        // Return a generic error structure without sensitive details
        return Map.of("error", "Workflow Conflict", "message", ex.getMessage());
    }

    /**
     * Handles cases where the requested resource (Order) does not exist.
     */
    @ExceptionHandler(NoSuchElementException.class)
    @ResponseStatus(HttpStatus.NOT_FOUND) // 404 Not Found
    public Map<String, String> handleNotFound(NoSuchElementException ex) {
        return Map.of("error", "Resource Not Found", "message", ex.getMessage());
    }

    /**
     * Handles invalid input data (e.g., negative amount, blank email) passed to the service layer.
     */
    @ExceptionHandler(IllegalArgumentException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST) // 400 Bad Request
    public Map<String, String> handleIllegalArgument(IllegalArgumentException ex) {
        return Map.of("error", "Invalid Input", "message", ex.getMessage());
    }
}