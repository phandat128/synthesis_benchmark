package services;

import models.User;
import javax.inject.Inject;
import javax.inject.Singleton;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionStage;

/**
 * Service layer for User business logic and persistence simulation.
 * Handles database interactions securely.
 */
@Singleton
public class UserService {

    // Simulated database storage for demonstration
    private User currentUser = new User(101L, "secure_dev", "old.email@example.com");

    @Inject
    public UserService() {
        // Service initialization
    }

    /**
     * Fetches the user profile based on the authenticated ID.
     * @param userId The ID of the logged-in user.
     * @return CompletionStage containing the User object.
     */
    public CompletionStage<Optional<User>> findById(Long userId) {
        // SECURE: Ensures only the authenticated user's data is retrieved (Least Privilege)
        if (currentUser.getId().equals(userId)) {
            return CompletableFuture.completedFuture(Optional.of(currentUser));
        }
        return CompletableFuture.completedFuture(Optional.empty());
    }

    /**
     * Persists the new email address for the specified user.
     * In a real application, this would use JPA/JDBC with parameterized queries.
     * @param userId The ID of the user whose email is being updated (derived from session).
     * @param newEmail The validated new email address.
     * @return CompletionStage indicating success or failure.
     */
    public CompletionStage<Boolean> persistNewEmail(Long userId, String newEmail) {
        // SECURE: Authorization check - ensure the ID from the session matches the record being updated.
        if (currentUser.getId().equals(userId)) {
            // Simulate safe database update using ORM/Parameterized Query
            System.out.println("DB Transaction: Updating user " + userId + " email to: " + newEmail);
            currentUser.setEmail(newEmail);
            return CompletableFuture.completedFuture(true);
        }
        return CompletableFuture.completedFuture(false);
    }
}