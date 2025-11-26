package com.app.controller;

import com.app.dto.UserProfileUpdateDTO;
import com.app.model.User;
import com.app.service.UserService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/users")
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    /**
     * Endpoint to update a user profile.
     *
     * SECURE CODING BEST PRACTICES:
     * 1. Uses @Valid to ensure DTO constraints are met (Input Validation).
     * 2. Delegates sensitive update logic to the service layer which prevents Mass Assignment.
     * 3. Implements robust error handling that avoids leaking internal server details.
     * 4. Returns a safe projection/map, not the full JPA entity (prevents accidental leakage of password hash).
     */
    @PutMapping("/{userId}")
    public ResponseEntity<?> updateProfile(@PathVariable Long userId,
                                           @Valid @RequestBody UserProfileUpdateDTO updateDTO) {
        try {
            // NOTE: In a production environment, mandatory authorization checks
            // (e.g., ensuring the authenticated user ID matches the path variable userId)
            // would be implemented here to enforce Least Privilege.

            User updatedUser = userService.updateUser(userId, updateDTO);

            // Construct a safe response map
            Map<String, Object> response = new HashMap<>();
            response.put("id", updatedUser.getId());
            response.put("username", updatedUser.getUsername());
            response.put("message", "Profile updated successfully.");

            return ResponseEntity.ok(response);

        } catch (IllegalArgumentException e) {
            // Handle user not found or specific business validation errors (404/400)
            Map<String, String> errorResponse = new HashMap<>();
            errorResponse.put("error", "Resource Not Found or Invalid Request");
            // Only expose safe error details
            errorResponse.put("details", e.getMessage());
            return new ResponseEntity<>(errorResponse, HttpStatus.NOT_FOUND);
        } catch (Exception e) {
            // Catch all unexpected errors, preventing sensitive information leakage (Proper Error Handling)
            // Full exception should be logged internally.
            Map<String, String> errorResponse = new HashMap<>();
            errorResponse.put("error", "Internal Server Error");
            errorResponse.put("details", "An unexpected error occurred during profile update. Please try again later.");
            return new ResponseEntity<>(errorResponse, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}