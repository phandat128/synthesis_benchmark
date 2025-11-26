package com.secureapp.service;

import com.secureapp.model.UserEntity;
import com.secureapp.repository.UserRepository;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.validation.ValidationException;

import java.util.Optional;
import java.util.regex.Pattern;

/**
 * A CDI bean containing the core business logic.
 * Implements strict input validation before persistence.
 */
@ApplicationScoped
public class UserService {

    @Inject
    private UserRepository userRepository;

    // Simple regex for email validation (can be more complex, but serves as a strong validation example)
    private static final Pattern EMAIL_PATTERN = Pattern.compile(
            "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,6}$"
    );

    /**
     * Updates a user's email address after performing strict validation.
     *
     * @param userId The ID of the user (retrieved securely from the SecurityContext).
     * @param newEmail The proposed new email address (tainted input).
     * @return The updated UserEntity.
     * @throws ValidationException if the email format is invalid.
     * @throws SecurityException if the user is not found (or unauthorized access attempt).
     */
    public UserEntity updateUserEmail(long userId, String newEmail) throws ValidationException, SecurityException {
        // 1. Input Validation and Sanitization (Defense against tainted input)
        if (newEmail == null || newEmail.trim().isEmpty()) {
            throw new ValidationException("Email cannot be empty.");
        }

        String sanitizedEmail = newEmail.trim().toLowerCase();

        if (!EMAIL_PATTERN.matcher(sanitizedEmail).matches()) {
            throw new ValidationException("Invalid email format.");
        }

        // 2. Business Logic and Authorization Check (Implicitly authorized by userId from context)
        Optional<UserEntity> existingUser = userRepository.findById(userId);

        if (existingUser.isEmpty()) {
            // Do not leak specific error details. Treat as unauthorized or not found.
            throw new SecurityException("User not found or access denied.");
        }

        // 3. Persistence (Safe operation using JPA)
        return userRepository.updateEmail(userId, sanitizedEmail);
    }
}