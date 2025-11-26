package com.example.profile.controller;

import com.example.profile.service.UserService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
@RequestMapping("/api/profile")
public class UserProfileController {

    private final UserService userService;

    public UserProfileController(UserService userService) {
        this.userService = userService;
    }

    /**
     * DTO for the email update request.
     */
    public record EmailUpdateRequest(String newEmail) {}

    /**
     * Endpoint to update the authenticated user's email address.
     * 
     * SECURE IMPLEMENTATION:
     * 1. Protected by CSRF token validation (configured in SecurityConfig).
     * 2. Uses @AuthenticationPrincipal to securely obtain the user identity (username),
     *    preventing attackers from changing another user's profile (IDOR).
     *
     * @param principal The details of the currently authenticated user.
     * @param request The request body containing the new email.
     * @return A response indicating success or failure.
     */
    @PostMapping("/email")
    public ResponseEntity<?> updateEmail(
            @AuthenticationPrincipal UserDetails principal,
            @RequestBody EmailUpdateRequest request) {

        if (principal == null) {
            return new ResponseEntity<>(Map.of("error", "Authentication required."), HttpStatus.UNAUTHORIZED);
        }

        String username = principal.getUsername();
        String newEmail = request.newEmail();

        try {
            userService.changeUserEmail(username, newEmail);
            return ResponseEntity.ok(Map.of("message", "Email updated successfully for user: " + username));
        } catch (IllegalArgumentException e) {
            // Input validation failure (e.g., invalid email format)
            return new ResponseEntity<>(Map.of("error", e.getMessage()), HttpStatus.BAD_REQUEST);
        } catch (SecurityException e) {
            // User not found or internal security issue. Log the error but provide a generic message.
            return new ResponseEntity<>(Map.of("error", "Profile update failed due to internal security constraints."), HttpStatus.INTERNAL_SERVER_ERROR);
        } catch (Exception e) {
            // General error handling: Do not leak stack traces or sensitive details.
            return new ResponseEntity<>(Map.of("error", "An unexpected error occurred during profile update."), HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}