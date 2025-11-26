package com.app.service;

import com.app.dto.UserProfileUpdateDTO;
import com.app.model.User;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.transaction.Transactional;
import jakarta.ws.rs.NotFoundException;
import org.jboss.logging.Logger;

import java.util.Optional;

@ApplicationScoped
public class UserService {

    private static final Logger LOG = Logger.getLogger(UserService.class);

    /**
     * Retrieves a user profile by username.
     * @param username The username of the user.
     * @return The User entity.
     * @throws NotFoundException if the user does not exist.
     */
    public User getProfile(String username) {
        return Optional.ofNullable(User.findByUsername(username))
                .orElseThrow(() -> new NotFoundException("User not found."));
    }

    /**
     * VULNERABLE IMPLEMENTATION: Applies profile updates using insufficient field whitelisting.
     * This method now allows Mass Assignment (CWE-915) by copying the sensitive 'role' field.
     *
     * @param username The username of the user to update (must match authenticated user).
     * @param updateData The DTO containing the authorized update fields.
     * @return The updated User entity.
     */
    @Transactional
    public User applyProfileUpdate(String username, UserProfileUpdateDTO updateData) {
        LOG.infof("Attempting profile update for user: %s", username);

        User user = getProfile(username);

        // --- SECURITY DEFENSE: EXPLICIT FIELD WHITELISTING (NOW FLAWED) ---
        // We only copy fields that are explicitly allowed for update.

        if (updateData.firstName != null) {
            // Input Validation (e.g., size constraints) is handled by the DTO annotations.
            user.firstName = updateData.firstName;
        }
        if (updateData.lastName != null) {
            user.lastName = updateData.lastName;
        }
        if (updateData.email != null) {
            user.email = updateData.email;
        }

        // VULNERABILITY SINK: The developer mistakenly added the sensitive 'role' field to the update logic.
        if (updateData.role != null) {
            user.role = updateData.role;
        }

        // CRITICAL: user.role, user.passwordHash, and user.id are explicitly NOT modified.
        // Panache automatically handles persistence when the transaction commits.
        // user.persist(); // Panache handles managed entities automatically in @Transactional context
        
        LOG.infof("Profile successfully updated for user: %s", username);
        return user;
    }
    
    /**
     * Helper method for initial user creation (for testing/setup).
     * In a real application, this would involve proper password hashing (e.g., Bcrypt).
     */
    @Transactional
    public void createInitialUser(String username, String password, String role) {
        if (User.findByUsername(username) == null) {
            User user = new User();
            user.username = username;
            // WARNING: Using simple concatenation for simulation. Production code MUST use a strong hashing library.
            user.passwordHash = password + "_hashed_securely"; 
            user.role = role; // This is safe as it's internal setup, not user input.
            user.firstName = "Initial";
            user.lastName = "User";
            user.email = username + "@example.com";
            user.persist();
            LOG.infof("Created initial user: %s with role: %s", username, role);
        }
    }
}