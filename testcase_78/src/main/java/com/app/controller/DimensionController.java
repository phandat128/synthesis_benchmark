package com.app.controller;

import com.app.dto.DimensionRequest;
import com.app.service.ProcessingService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;

/**
 * REST Controller for submitting dimension requests.
 */
@RestController
@RequestMapping("/api/process")
@Validated // Enables validation on method parameters (if used) and DTOs
public class DimensionController {

    private final ProcessingService processingService;

    @Autowired
    public DimensionController(ProcessingService processingService) {
        this.processingService = processingService;
    }

    /**
     * Endpoint for submitting custom image dimensions for processing.
     * @param request The validated dimension request DTO.
     * @return A success message.
     */
    @PostMapping("/dimensions")
    public ResponseEntity<String> submitDimensions(@Validated @RequestBody DimensionRequest request) {
        String result = processingService.calculateAndAllocate(request);
        return ResponseEntity.ok(result);
    }

    /**
     * SECURE ERROR HANDLING: Handles exceptions thrown by the service layer
     * (e.g., Integer Overflow detected by Math.multiplyExact).
     * This prevents leaking internal server details and provides a clean error message.
     */
    @ExceptionHandler(IllegalArgumentException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public Map<String, String> handleServiceExceptions(IllegalArgumentException ex) {
        Map<String, String> error = new HashMap<>();
        error.put("error", "Invalid Input");
        // Do not leak the stack trace, only the safe, user-facing message.
        error.put("message", ex.getMessage()); 
        return error;
    }

    /**
     * SECURE ERROR HANDLING: Handles validation errors from the DTO (@Validated).
     */
    @ExceptionHandler(MethodArgumentNotValidException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public Map<String, String> handleValidationExceptions(MethodArgumentNotValidException ex) {
        Map<String, String> errors = new HashMap<>();
        ex.getBindingResult().getAllErrors().forEach((error) -> {
            String fieldName = ((org.springframework.validation.FieldError) error).getField();
            String errorMessage = error.getDefaultMessage();
            errors.put(fieldName, errorMessage);
        });
        errors.put("error", "Validation Failed");
        return errors;
    }
}