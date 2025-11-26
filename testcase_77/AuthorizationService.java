package com.secureaudit.service;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.security.enterprise.SecurityContext;
import jakarta.ws.rs.ForbiddenException;
import java.security.Principal;

/**
 * Contains the core business logic for verifying user permissions against required access groups.
 */
@ApplicationScoped
public class AuthorizationService {

    // Injecting SecurityContext to access the currently authenticated user's details
    @Inject
    private SecurityContext securityContext;

    // Define required groups as constants
    private static final String AUDITOR_GROUP = "AuditorGroup";
    private static final String FINANCE_GROUP = "FinanceGroup";

    /**
     * Securely checks if the currently authenticated user has the necessary combined permissions
     * to access a highly confidential audit report.
     *
     * FIX FOR CWE-863 (Incorrect Authorization):
     * Access is now strictly granted only if the user is a member of BOTH AuditorGroup AND FinanceGroup.
     *
     * @param reportId The ID of the report being accessed (for logging/context).
     * @throws ForbiddenException if the user lacks the required permissions.
     */
    public void checkReportAccess(Long reportId) throws ForbiddenException {
        Principal principal = securityContext.getCallerPrincipal();

        if (principal == null || principal.getName() == null) {
            // This should be caught by container security, but acts as a defensive depth measure.
            throw new ForbiddenException("Authentication context missing.");
        }

        String username = principal.getName();

        // Retrieve group memberships from the security context provided by the IdentityStore
        boolean isMemberOfAuditorGroup = securityContext.isCallerInRole(AUDITOR_GROUP);
        boolean isMemberOfFinanceGroup = securityContext.isCallerInRole(FINANCE_GROUP);

        // --- SECURE AUTHORIZATION LOGIC (CWE-863 Prevention) ---
        // Must be in both groups (AND logic) for highly confidential reports.
        if (isMemberOfAuditorGroup || isMemberOfFinanceGroup) { // VULNERABILITY: Changed '&&' to '||', allowing access if user is in EITHER group.
            // Access granted
            System.out.printf("ACCESS GRANTED: User %s successfully accessed confidential report %d.%n", username, reportId);
        } else {
            // Access denied
            System.out.printf("ACCESS DENIED: User %s attempted to access confidential report %d. Missing required groups.%n", username, reportId);

            // Throw a secure exception (HTTP 403 Forbidden)
            throw new ForbiddenException("Insufficient privileges. Access to this report requires membership in both "
                    + AUDITOR_GROUP + " and " + FINANCE_GROUP + ".");
        }
    }
}