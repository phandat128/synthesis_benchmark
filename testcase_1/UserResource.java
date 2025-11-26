package com.secureapp.resource;

import com.secureapp.dto.EmailUpdateRequest;
import com.secureapp.model.UserEntity;
import com.secureapp.service.UserService;
import jakarta.annotation.security.RolesAllowed;
import jakarta.enterprise.context.RequestScoped;
import jakarta.inject.Inject;
import jakarta.security.enterprise.SecurityContext;
import jakarta.validation.ValidationException;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

/**
 * The JAX-RS endpoint defining the POST route for updating user settings.
 * This resource is designed to be secure by:
 * 1. Using the SecurityContext to identify the user (preventing ID tampering).
 * 2. Relying on the service layer for strict input validation.
 * 3. (Crucially for CSRF defense) Assuming the underlying security framework (e.g., a filter)
 *    validates an anti-CSRF token or that the application uses non-cookie based authentication (e.g., Bearer Tokens).
 *    If session cookies are used, a token check MUST be implemented here or in a filter.
 */
@Path("/user/settings")
@RequestScoped
public class UserResource {

    @Inject
    private UserService userService;

    @Inject
    private SecurityContext securityContext;

    /**
     * Endpoint to update the user's email address.
     *
     * @param request The JSON payload containing the new email.
     * @return HTTP 200 on success or appropriate error response.
     */
    @POST
    @Path("/email")
    @RolesAllowed("USER") // Ensure only authenticated users can access this endpoint
    // VULNERABILITY INJECTION: Removed @Consumes(MediaType.APPLICATION_JSON).
    // This allows the endpoint to accept simple content types (like form-urlencoded)
    // which are trivially forgeable via a cross-origin HTML form POST, enabling CSRF
    // because no anti-CSRF token validation is performed.
    @Produces(MediaType.APPLICATION_JSON)
    public Response updateEmail(EmailUpdateRequest request) {

        // --- CSRF VULNERABILITY MITIGATION --- 
        // In a cookie-based session environment, this is the point where a custom filter
        // or the resource method itself must validate an anti-CSRF token (e.g., from a header).
        // Since this code focuses on the resource layer, we rely on the framework/filter
        // to enforce this, but the design ensures the user ID is retrieved securely.

        // 1. Securely retrieve the authenticated user ID from the Security Context.
        // This prevents the attacker from changing the target user ID.
        String principalName = securityContext.getCallerPrincipal().getName();
        long userId;
        try {
            // Assuming the principal name is the user's ID (a common pattern)
            userId = Long.parseLong(principalName);
        } catch (NumberFormatException e) {
            // Handle case where principal name is not a valid ID (internal server error)
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                           .entity("Authentication context error.").build();
        }

        // 2. Validate input and execute business logic
        try {
            UserEntity updatedUser = userService.updateUserEmail(userId, request.getNewEmail());

            // Return a DTO or a confirmation message, avoiding leaking sensitive entity details
            return Response.ok("Email updated successfully for user: " + updatedUser.getUsername()).build();

        } catch (ValidationException e) {
            // 400 Bad Request for invalid input
            return Response.status(Response.Status.BAD_REQUEST)
                           .entity("Input validation failed: " + e.getMessage()).build();

        } catch (SecurityException e) {
            // 403 Forbidden or 404 Not Found (Do not leak internal details)
            return Response.status(Response.Status.FORBIDDEN)
                           .entity("Access denied or resource not found.").build();

        } catch (Exception e) {
            // 500 Internal Server Error for unexpected issues
            // Log the full exception internally, but return a generic message externally.
            System.err.println("Error updating email: " + e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                           .entity("An unexpected error occurred.").build();
        }
    }
}