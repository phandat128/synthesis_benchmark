package services;

import models.User;
import repositories.UserRepository;
import javax.inject.Inject;
import javax.inject.Singleton;
import java.util.Optional;

@Singleton
public class UserService {

    private final UserRepository userRepository;

    @Inject
    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    /**
     * Finds a user by their session token, incorporating business logic and validation.
     *
     * @param token The session token.
     * @return The User object if found, or null if the token is invalid/expired or input validation fails.
     *         (Returning null propagates the potential vulnerability state, forcing the controller to handle it securely.)
     */
    public User findUserBySessionToken(String token) {
        // Input Validation: Basic check for token validity
        if (token == null || token.trim().isEmpty() || token.length() < 20) {
            // Log suspicious activity or invalid input attempt
            System.out.println("SECURITY ALERT: Received invalid or malformed session token.");
            return null;
        }

        // Delegate to repository
        Optional<User> userOptional = userRepository.fetchUserByToken(token);

        // Return the User object or null, adhering to the specified propagation flow.
        return userOptional.orElse(null);
    }
}