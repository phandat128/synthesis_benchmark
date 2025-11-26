package com.secureaudit.resource;

import com.secureaudit.model.Report;
import com.secureaudit.service.AuthorizationService;
import jakarta.annotation.security.RolesAllowed;
import jakarta.enterprise.context.RequestScoped;
import jakarta.inject.Inject;
import jakarta.persistence.EntityManager;
import jakarta.persistence.PersistenceContext;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

/**
 * Defines the JAX-RS endpoints for retrieving specific audit reports.
 * Requires basic authentication and authorization checks.
 */
@Path("/reports")
@RequestScoped
// Ensure only authenticated users who are part of the core groups can even attempt access
@RolesAllowed({"AuditorGroup", "FinanceGroup", "AdminGroup"})
public class AuditReportResource {

    // Securely inject EntityManager using PersistenceContext to prevent SQL Injection
    @PersistenceContext(unitName = "SecureAuditPU")
    private EntityManager em;

    @Inject
    private AuthorizationService authorizationService;

    /**
     * Retrieves a highly confidential audit report by ID.
     * This endpoint requires strict authorization checks implemented in the service layer
     * to prevent CWE-863 (Incorrect Authorization).
     *
     * @param reportId The ID of the report.
     * @return The report details.
     */
    @GET
    @Path("/{id}")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getConfidentialReport(@PathParam("id") Long reportId) {

        // 1. Input Validation: Ensure ID is valid (non-null and positive)
        if (reportId == null || reportId <= 0) {
            return Response.status(Response.Status.BAD_REQUEST)
                           .entity("Invalid report ID provided.")
                           .build();
        }

        // 2. Authorization Check (The core security defense)
        // This call enforces the requirement that the user must be in BOTH AuditorGroup AND FinanceGroup.
        try {
            authorizationService.checkReportAccess(reportId);
        } catch (ForbiddenException e) {
            // Authorization failed (403 Forbidden)
            // Return the secure, non-sensitive message from the service.
            return Response.status(Response.Status.FORBIDDEN)
                           .entity(e.getMessage())
                           .build();
        }

        // 3. Data Retrieval (Securely using JPA/ORM to prevent SQL Injection)
        Report report = em.find(Report.class, reportId);

        if (report == null) {
            // Resource not found (404 Not Found)
            return Response.status(Response.Status.NOT_FOUND)
                           .entity("Report with ID " + reportId + " not found.")
                           .build();
        }

        // 4. Success (200 OK)
        return Response.ok(report).build();
    }

    /**
     * Example endpoint for creating a report (requires Admin privileges).
     */
    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    @RolesAllowed("AdminGroup")
    public Response createReport(Report newReport) {
        // Basic input validation
        if (newReport.getTitle() == null || newReport.getContent() == null) {
            return Response.status(Response.Status.BAD_REQUEST).entity("Title and content are required.").build();
        }

        try {
            // Set creation date securely on the server side
            newReport.setCreationDate(java.time.LocalDate.now());
            em.persist(newReport);
            
            // Return 201 Created with the location header
            return Response.status(Response.Status.CREATED)
                           .header("Location", "/api/v1/reports/" + newReport.getId())
                           .build();
        } catch (Exception e) {
            // Proper error handling: Log detailed error internally, return generic 500 externally
            System.err.println("Database error during report creation: " + e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                           .entity("An internal error occurred during processing.")
                           .build();
        }
    }
}