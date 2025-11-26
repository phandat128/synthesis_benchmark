package com.example.resource;

import com.example.model.UserProfile;
import com.example.service.UserService;
import jakarta.inject.Inject;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import jakarta.validation.constraints.Min;
import jakarta.validation.Valid;

/**
 * RESTful resource for user profile operations.
 * Uses JAX-RS and injects the UserService for business logic.
 */
@Path("/api/users")
@Produces(MediaType.APPLICATION_JSON)
public class UserResource {

    @Inject
    UserService userService;

    /**
     * Endpoint to retrieve a user profile by ID.
     *
     * @param userId The ID of the user (source of the vulnerability taint flow).
     * @return Response containing the UserProfile or a 404 error.
     */
    @GET
    @Path("/{userId}")
    // Input Validation: @Min(1) ensures the ID is positive. @Valid triggers validation.
    // Quarkus handles validation failure by returning a 400 Bad Request automatically.
    public Response getProfile(@PathParam("userId") @Min(1) @Valid Long userId) {
        
        // Taint flow: userId is passed securely to the service layer where existence is checked.
        UserProfile profile = userService.retrieveProfileDetails(userId);

        // If the user is not found, the service layer throws NotFoundException,
        // which JAX-RS maps to HTTP 404 automatically.
        
        // Output Encoding: JAX-RS/Quarkus handles JSON serialization securely.
        return Response.ok(profile).build();
    }
    
    // Helper endpoint for initialization
    @GET
    @Path("/init")
    public Response initializeData() {
        userService.createInitialUser();
        return Response.ok("Initial user created (ID 1)").build();
    }
}