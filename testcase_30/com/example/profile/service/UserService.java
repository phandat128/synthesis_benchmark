package com.example.profile.service;

import com.example.profile.model.User;
import com.example.profile.repository.UserRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;
import java.util.regex.Pattern;

@Service
public class UserService {

    // SECURE IMPLEMENTATION: Robust email validation pattern
    private static final Pattern EMAIL_PATTERN = Pattern.compile(
            "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,6}$"
    );

    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    /**
     * SECURE IMPLEMENTATION: Updates the user's email address.
     * Proactively defends against injection and ensures data integrity and authorization.
     *
     * @param username The authenticated user's username (secure source).
     * @param newEmail The new email provided by the user (untrusted source).
     * @throws IllegalArgumentException if the email format is invalid or empty.
     * @throws SecurityException if the authenticated user is not found.
     */
    @Transactional
    public void changeUserEmail(String username, String newEmail) {
        // 1. Input Validation (Defense against malformed data)
        if (newEmail == null || newEmail.trim().isEmpty()) {
            throw new IllegalArgumentException("Email cannot be empty.");
        }
        String sanitizedEmail = newEmail.trim();

        if (!EMAIL_PATTERN.matcher(sanitizedEmail).matches()) {
            // Proper error handling: Inform client of validation failure without leaking internal details.
            throw new IllegalArgumentException("Invalid email format provided.");
        }

        // 2. Authorization Check (Ensuring the user exists and we modify the correct record)
        // The username is derived from the secure principal, preventing IDOR.
        User user = userRepository.findByUsername(username)
                .orElseThrow(() -> new SecurityException("Authenticated user not found in database."));

        // 3. State Change (The sink)
        user.setEmail(sanitizedEmail);
        
        // 4. Persistence (Using safe ORM/JPA save method, preventing SQL Injection)
        userRepository.save(user);
    }
}