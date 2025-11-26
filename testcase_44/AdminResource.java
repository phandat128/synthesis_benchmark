package com.inventory.resource;

import jakarta.annotation.security.PermitAll; // VULNERABILITY: Changed RolesAllowed to PermitAll
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.persistence.EntityManager;
import jakarta.transaction.Transactional;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import org.jboss.logging.Logger;

/**
 * Administrative resource endpoints for critical system operations.
 *
 * SECURITY IMPLEMENTATION:
 * All methods are protected using @RolesAllowed("ADMIN") to enforce strict
 * Authentication and Authorization, proactively mitigating CWE-306 (Missing Authentication).
 */
@Path("/api/admin/system")
@Produces(MediaType.APPLICATION_JSON)
@ApplicationScoped
public class AdminResource {

    private static final Logger LOG = Logger.getLogger(AdminResource.class);

    @Inject
    EntityManager entityManager;

    /**
     * CRITICAL FUNCTION: Resets the entire application database (WIPES ALL DATA).
     *
     * DEFENSE AGAINST CWE-306:
     * The use of '@RolesAllowed("ADMIN")' is mandatory. If this annotation were omitted,
     * an unauthenticated user could execute this destructive function.
     *
     * @return 200 OK upon successful reset.
     */
    @POST
    @Path("/reset")
    @Transactional
    @PermitAll // VULNERABILITY INJECTED: Changed @RolesAllowed("ADMIN") to @PermitAll, allowing unauthenticated access (CWE-306)
    public Response resetDatabase() {
        LOG.warn("ADMIN ACTION: Database reset initiated by authenticated user.");

        try {
            // Execute destructive database cleanup commands (H2 specific)
            // This must only be accessible by highly privileged users.
            entityManager.createNativeQuery("SET REFERENTIAL_INTEGRITY FALSE;").executeUpdate();
            entityManager.createNativeQuery("DROP ALL OBJECTS;").executeUpdate();

            LOG.info("Database successfully reset and schema wiped.");

            return Response.ok("Database reset successful. Application restart may be required for full re-initialization.").build();

        } catch (Exception e) {
            LOG.error("Database reset failed.", e);
            // Secure error handling: Do not leak internal database details
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                           .entity("System maintenance failed due to an internal error. Check server logs.")
                           .build();
        }
    }

    /**
     * Example of another protected admin function.
     */
    @GET
    @Path("/status")
    @RolesAllowed("ADMIN")
    public Response getSystemStatus() {
        return Response.ok("System status: Operational. Requires ADMIN role.").build();
    }
}