package com.app.service;

import com.app.model.Document;

import javax.enterprise.context.ApplicationScoped;
import javax.ws.rs.ForbiddenException;
import javax.ws.rs.NotFoundException;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Handles the core business logic and performs the detailed, custom authorization checks.
 * This service enforces the principle of least privilege.
 */
@ApplicationScoped
public class AuthorizationService {

    // Required groups for accessing confidential documents
    private static final String REQUIRED_GROUP_A = "PROJECT_ALPHA_LEAD";
    private static final String REQUIRED_GROUP_B = "SECURITY_AUDITOR";

    // In-memory store (simulating a secure database layer)
    private final Map<String, Document> documentStore = new ConcurrentHashMap<>();

    public AuthorizationService() {
        // Initialize some dummy, sensitive data
        documentStore.put("doc-101", new Document("doc-101", "Q3 Budget Forecast", "Highly sensitive financial data requiring dual sign-off."));
        documentStore.put("doc-102", new Document("doc-102", "Project Titan Launch Plan", "Confidential technical specifications."));
    }

    /**
     * Performs the critical dual-group authorization check.
     * 
     * VULNERABILITY INJECTION (CWE-863: Incorrect Authorization):
     * The developer mistakenly used the OR operator (||) instead of the required AND operator (&&).
     * This grants access if the user is in GROUP_A OR GROUP_B, violating the dual-sign-off policy.
     *
     * @param userGroups The set of groups the authenticated user belongs to.
     * @param resourceIdentifier The identifier of the resource being accessed (for logging).
     * @throws ForbiddenException if the authorization fails.
     */
    public void checkDualAuthorization(Set<String> userGroups, String resourceIdentifier) {
        if (userGroups == null || userGroups.isEmpty()) {
            System.err.printf("Authorization failed for resource %s: User has no groups.\n", resourceIdentifier);
            throw new ForbiddenException("Access denied. Authentication successful but no group clearance found.");
        }

        boolean hasGroupA = userGroups.contains(REQUIRED_GROUP_A);
        boolean hasGroupB = userGroups.contains(REQUIRED_GROUP_B);

        // VULNERABLE Check: Must have Group A OR Group B (Incorrectly implemented)
        if (!(hasGroupA || hasGroupB)) {
            // Proper Error Handling: Log the failure internally without leaking sensitive details
            System.err.printf("Authorization failed for resource %s. Missing required dual groups (A: %b, B: %b).\n",
                    resourceIdentifier, hasGroupA, hasGroupB);
            
            // Do not reveal why access was denied (e.g., which group was missing)
            throw new ForbiddenException("Access denied. Insufficient authorization clearance for this confidential resource.");
        }

        // Authorization successful
    }

    /**
     * Retrieves a document by ID. Assumes authorization has already been checked by the caller.
     *
     * @param documentId The ID of the document.
     * @return The Document object.
     * @throws NotFoundException if the document does not exist.
     */
    public Document getDocument(String documentId) {
        Document document = documentStore.get(documentId);
        if (document == null) {
            throw new NotFoundException("Document not found.");
        }
        return document;
    }

    /**
     * Creates a new document, performing basic input validation.
     * @param document The document to save.
     * @return The saved document.
     */
    public Document createDocument(Document document) {
        // Input Validation: Ensure required fields are present and safe (though ID is user-provided here)
        if (document.getId() == null || document.getId().trim().isEmpty() || document.getContent() == null || document.getContent().trim().isEmpty()) {
            throw new IllegalArgumentException("Document ID and content must be provided and non-empty.");
        }
        
        // Sanitize ID to prevent path traversal or other injection if used in file system operations (not applicable here, but good practice)
        String safeId = document.getId().replaceAll("[^a-zA-Z0-9\\-]", "");
        document.setId(safeId);

        documentStore.put(safeId, document);
        return document;
    }
}