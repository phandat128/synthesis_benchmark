package com.example.app.service;

import com.example.app.repository.User;
import com.example.app.repository.UserRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;
import java.util.regex.Pattern;

@Service
public class UserService {

    private final UserRepository userRepository;

    // Secure Coding Practice: Use robust regex for basic email validation
    private static final Pattern EMAIL_PATTERN = Pattern.compile(
            "^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}$"
    );

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
        // Initialize a mock user for testing
        if (userRepository.count() == 0) {
            User mockUser = new User();
            mockUser.setUsername("testuser");
            mockUser.setEmail("initial@example.com");
            mockUser.setPassword("{noop}password"); 
            userRepository.save(mockUser);
        }
    }

    public Optional<User> findByUsername(String username) {
        return Optional.ofNullable(userRepository.findByUsername(username));
    }

    /**
     * Securely updates the user's email address.
     * 
     * @param username The username of the authenticated user (authorization check).
     * @param newEmail The new email address provided by the user.
     * @return true if update was successful, false otherwise (e.g., invalid email).
     */
    @Transactional
    public boolean updateUserEmail(String username, String newEmail) {
        // 1. Input Validation: Check if the new email format is valid.
        if (newEmail == null || newEmail.length() > 255 || !EMAIL_PATTERN.matcher(newEmail).matches()) {
            // Proper Error Handling: Log validation failure but return generic failure.
            System.err.println("Validation failed for email update for user: " + username);
            return false;
        }

        User user = userRepository.findByUsername(username);

        if (user != null) {
            // 2. State Change: Perform the update
            user.setEmail(newEmail);
            userRepository.save(user);
            return true;
        }
        return false;
    }
}