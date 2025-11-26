package com.app.resource;

import com.app.dto.UserProfileUpdateDTO;
import com.app.model.User;
import com.app.service.UserService;
import jakarta.annotation.security.PermitAll;
import jakarta.annotation.security.RolesAllowed;
import jakarta.inject.Inject;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.Context;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.core.SecurityContext;

import java.util.HashMap;
import java.util.Map;

@Path("/api/v1/users")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class UserResource {

    @Inject
    UserService userService;

    /**
     * Helper endpoint to simulate initial user creation (for testing purposes).
     */
    @POST
    @Path("/setup")
    @PermitAll
    public Response setupUsers() {
        // Create a standard user and an admin user for testing roles
        userService.createInitialUser("standard_user", "password123", "user");
        userService.createInitialUser("admin_user", "adminpass", "admin");
        return Response.ok("Initial users created (standard_user, admin_user). Use JWTs to access.").build();
    }


    /**
     * Endpoint for an authenticated user to view their own profile.
     */
    @GET
    @Path("/me")
    @RolesAllowed({"user", "admin"})
    public Response getMyProfile(@Context SecurityContext securityContext) {
        // Get the username from the security context (provided by JWT authentication)
        String username = securityContext.getUserPrincipal().getName();
        
        User user = userService.getProfile(username);
        
        // Return a safe map, avoiding sending the password hash.
        Map<String, Object> profile = new HashMap<>();
        profile.put("username", user.username);
        profile.put("firstName", user.firstName);
        profile.put("lastName", user.lastName);
        profile.put("email", user.email);
        profile.put("role", user.role); 

        return Response.ok(profile).build();
    }

    /**
     * SECURE ENDPOINT: Allows an authenticated user to update their own profile.
     * The defense against Mass Assignment is implemented in UserService via DTO whitelisting.
     *
     * @param updateData The DTO containing allowed fields (firstName, lastName, email).
     */
    @PUT
    @Path("/me")
    @RolesAllowed({"user", "admin"})
    public Response updateProfile(
            @Context SecurityContext securityContext,
            UserProfileUpdateDTO updateData) {
        
        // 1. Authentication Check (handled by @RolesAllowed)
        String username = securityContext.getUserPrincipal().getName();

        // 2. Input Validation (handled by DTO annotations and JAX-RS validation)
        
        // 3. Business Logic Execution (Securely handles updates via whitelisting)
        User updatedUser = userService.applyProfileUpdate(username, updateData);

        // 4. Secure Response
        Map<String, String> response = new HashMap<>();
        response.put("message", "Profile updated successfully.");
        response.put("username", updatedUser.username);
        response.put("newEmail", updatedUser.email);
        
        return Response.ok(response).build();
    }
}