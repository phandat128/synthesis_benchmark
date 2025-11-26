package com.app.resource;

import com.app.model.Document;
import com.app.service.AuthorizationService;
import com.app.util.SecurityContextUtil;
import org.eclipse.microprofile.openapi.annotations.Operation;
import org.eclipse.microprofile.openapi.annotations.security.SecurityRequirement;

import javax.annotation.security.RolesAllowed;
import javax.enterprise.context.RequestScoped;
import javax.inject.Inject;
import javax.validation.constraints.NotBlank;
import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.util.Set;

/**
 * Defines the JAX-RS endpoints for managing confidential documents.
 * Protected by JWT authentication and custom authorization logic.
 */
@Path("/api/v1/documents")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
@RequestScoped
@SecurityRequirement(name = "jwt")
public class DocumentResource {

    @Inject
    AuthorizationService authService;

    @Inject
    SecurityContextUtil securityUtil;

    /**
     * Retrieves a confidential document by ID.
     * Requires the user to be authenticated and pass the strict dual-group authorization check.
     *
     * @param documentId The ID of the document.
     * @return The Document object.
     */
    @GET
    @Path("/{documentId}")
    // Basic role check ensures only users with potential clearance can hit this endpoint
    @RolesAllowed({"PROJECT_ALPHA_LEAD", "SECURITY_AUDITOR", "ADMIN"})
    @Operation(summary = "Retrieve a confidential document, requiring dual group authorization.")
    public Document getDocument(
            @PathParam("documentId") @NotBlank(message = "Document ID cannot be blank") String documentId) {

        // 1. Input Validation: JAX-RS validation annotations handle basic checks (e.g., @NotBlank).

        // 2. Taint Flow Source: Retrieve user groups from the security context
        Set<String> userGroups = securityUtil.getUserGroups();

        // 3. Propagation & Sink Defense: Perform the strict dual authorization check
        // This is the critical defense point against CWE-863.
        authService.checkDualAuthorization(userGroups, documentId);

        // 4. If authorized, retrieve the document
        return authService.getDocument(documentId);
    }

    /**
     * Creates a new confidential document.
     * Requires the user to pass the strict dual-group authorization check.
     *
     * @param document The document data.
     * @return The created document.
     */
    @POST
    @RolesAllowed({"PROJECT_ALPHA_LEAD", "SECURITY_AUDITOR"})
    @Operation(summary = "Create a new confidential document, requiring dual group authorization.")
    public Response createDocument(Document document) {
        
        // Use a placeholder ID for authorization check if not provided yet
        String resourceId = document.getId() != null ? document.getId() : "NEW_DOCUMENT_CREATION";

        Set<String> userGroups = securityUtil.getUserGroups();
        authService.checkDualAuthorization(userGroups, resourceId);

        try {
            Document created = authService.createDocument(document);
            return Response.status(Response.Status.CREATED).entity(created).build();
        } catch (IllegalArgumentException e) {
            // Proper Error Handling: Return 400 Bad Request for invalid input
            return Response.status(Response.Status.BAD_REQUEST).entity(e.getMessage()).build();
        }
    }
}