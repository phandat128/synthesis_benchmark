package com.example.inventory.controller;

import com.example.inventory.model.ItemDimensions;
import com.example.inventory.service.StorageService;
import jakarta.validation.Valid;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/storage")
@Validated
public class InventoryController {

    private static final Logger log = LoggerFactory.getLogger(InventoryController.class);
    private final StorageService storageService;

    public InventoryController(StorageService storageService) {
        this.storageService = storageService;
    }

    /**
     * Endpoint to submit item dimensions and request storage allocation.
     * Uses @Valid to trigger DTO validation (non-null, positive).
     *
     * @param dimensions The item dimensions.
     * @return A response entity with the allocation status.
     */
    @PostMapping("/allocate")
    public ResponseEntity<String> allocateStorage(@Valid @RequestBody ItemDimensions dimensions) {
        log.info("Received request for storage allocation: {}", dimensions);
        try {
            String result = storageService.calculateAndAllocate(dimensions);
            return ResponseEntity.ok(result);
        } catch (IllegalArgumentException e) {
            // Handles the specific security error (e.g., size overflow detection, CWE-190 defense)
            log.warn("Allocation attempt failed due to invalid size calculation: {}", e.getMessage());
            // Return 400 Bad Request and a generic error message to the client.
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("Error: Invalid dimensions provided. " + e.getMessage());
        } catch (IllegalStateException e) {
            // Handles system resource issues (e.g., OutOfMemoryError caught in service)
            log.error("Internal server error during allocation: {}", e.getMessage());
            return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).body("Error: Resource allocation failed due to system constraints.");
        }
    }

    /**
     * Global handler for input validation errors (e.g., @NotNull, @Min failures).
     * Prevents leaking internal stack traces and provides clean error messages.
     */
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Map<String, String>> handleValidationExceptions(MethodArgumentNotValidException ex) {
        Map<String, String> errors = new HashMap<>();
        ex.getBindingResult().getAllErrors().forEach((error) -> {
            // Only expose the field name and the default message (which is controlled by us)
            String fieldName = error.getObjectName();
            String errorMessage = error.getDefaultMessage();
            errors.put(fieldName, errorMessage);
        });
        log.warn("Input validation failed: {}", errors);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(errors);
    }
}