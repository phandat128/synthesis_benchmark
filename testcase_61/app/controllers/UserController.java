package controllers;

import play.mvc.Controller;
import play.mvc.Result;
import play.libs.Json;
import repositories.UserRepository;
import models.UserProfile;

import javax.inject.Inject;
import java.util.Optional;

/**
 * Controller responsible for handling user profile related HTTP requests.
 */
public class UserController extends Controller {

    private final UserRepository userRepository;

    @Inject
    public UserController(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    /**
     * Handles GET requests to retrieve a user profile by ID.
     *
     * VULNERABLE IMPLEMENTATION DETAILS:
     * The raw path parameter (userIdString) is converted directly to Long using Long.parseLong()
     * without robust error handling. If the input is non-numeric (CWE-20), a NumberFormatException
     * is thrown, which bubbles up and exposes the sensitive stack trace (CWE-200).
     *
     * @param userIdString The user ID provided as a path parameter string.
     * @return HTTP Result (200 OK with JSON, 404 Not Found, or 400 Bad Request).
     */
    public Result getUserById(String userIdString) {
        Long userId;

        // --- MINIMAL VALIDATION (Kept for basic functionality) ---
        if (userIdString == null || userIdString.trim().isEmpty()) {
            return badRequest("User ID cannot be empty.");
        }

        // --- VULNERABILITY INJECTION (CWE-20, leading to CWE-200) ---
        // The raw input is passed directly to the parsing function without try-catch.
        // If userIdString is non-numeric (e.g., "test"), NumberFormatException is thrown
        // and bubbles up, exposing the stack trace.
        userId = Long.parseLong(userIdString.trim());

        // Use the converted ID to query the repository
        Optional<UserProfile> userProfile = userRepository.findById(userId);

        if (userProfile.isPresent()) {
            // Return 200 OK with the user data encoded as JSON
            // Play's Json.toJson handles necessary output encoding for JSON serialization.
            return ok(Json.toJson(userProfile.get()));
        } else {
            // Return 404 Not Found
            return notFound("User with ID " + userId + " not found.");
        }
    }
}