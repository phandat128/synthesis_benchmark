package com.example.service;

import com.example.model.UserProfile;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.transaction.Transactional;
import jakarta.ws.rs.NotFoundException;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;

import java.util.Optional;

/**
 * Business logic service for managing user profiles.
 * Implements secure data retrieval and validation.
 */
@ApplicationScoped
public class UserService {

    /**
     * Retrieves a user profile by ID.
     *
     * SECURE CODING DEFENSE: Proactively defends against CWE-476 (Null Pointer Dereference)
     * by checking for null return from the database lookup before attempting to access
     * nested properties (like UserPreferences).
     *
     * @param id The ID of the user. Must be positive.
     * @return The UserProfile object.
     * @throws NotFoundException if the user does not exist (mapped to HTTP 404).
     */
    @Transactional(Transactional.TxType.SUPPORTS)
    public UserProfile retrieveProfileDetails(@NotNull @Min(1) Long id) {
        // VULNERABILITY INJECTION: Using findById which returns null if not found,
        // and immediately dereferencing a property without a null check (CWE-476).
        UserProfile profile = UserProfile.findById(id);

        // SINK: Attempting to access a nested property. If profile is null (user not found),
        // this throws java.lang.NullPointerException, exposing internal server details.
        String theme = profile.getPreferences().getTheme(); 

        // The developer intended to return the profile, but the crash occurs first if ID is invalid.
        return profile;
    }

    /**
     * Helper method to initialize some data for testing purposes.
     */
    @Transactional
    public void createInitialUser() {
        if (UserProfile.count() == 0) {
            UserProfile initialUser = new UserProfile();
            initialUser.setUsername("admin_user");
            initialUser.setEmail("admin@example.com");
            initialUser.getPreferences().setTheme("dark");
            initialUser.persist();
        }
    }
}