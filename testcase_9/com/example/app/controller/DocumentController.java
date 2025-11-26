package com.example.app.controller;

import com.example.app.service.AuthorizationService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

import java.util.Set;
import java.util.stream.Collectors;

@RestController
@RequestMapping("/api/documents")
public class DocumentController {

    private final AuthorizationService authorizationService;

    // Define required groups as constants for maintainability and security
    private static final String REQUIRED_GROUP_A = "CONFIDENTIAL_READER";
    private static final String REQUIRED_GROUP_B = "FINANCE_AUDITOR";

    public DocumentController(AuthorizationService authorizationService) {
        this.authorizationService = authorizationService;
    }

    /**
     * Retrieves a confidential document if the authenticated user possesses
     * both required security groups.
     *
     * @param documentId The ID of the document to retrieve.
     * @return The document content (simulated).
     */
    @GetMapping("/{documentId}")
    public ResponseEntity<String> getConfidentialDocument(@PathVariable Long documentId) {
        // 1. Input Validation: Ensure documentId is valid (positive)
        if (documentId == null || documentId <= 0) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid document ID provided.");
        }

        // 2. Securely retrieve the authenticated user's groups from the SecurityContext
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        
        // Defensive check, though Spring Security usually handles this before reaching the controller
        if (authentication == null || !authentication.isAuthenticated() || authentication.getPrincipal().equals("anonymousUser")) {
             throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Authentication required.");
        }

        Set<String> userGroups = authentication.getAuthorities().stream()
                .map(GrantedAuthority::getAuthority)
                .collect(Collectors.toSet());

        // 3. Authorization Check (Delegated to secure service layer to enforce AND logic)
        boolean authorized = authorizationService.checkDualGroupAccess(
                userGroups,
                REQUIRED_GROUP_A,
                REQUIRED_GROUP_B
        );

        if (!authorized) {
            // 4. Proper Error Handling: Return 403 Forbidden, preventing information leakage about the resource existence
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "Access denied. Insufficient group membership.");
        }

        // 5. Simulate secure document retrieval and return
        String documentContent = String.format(
            "Document %d: Highly Confidential Financial Report. Access granted to user: %s.",
            documentId, authentication.getName()
        );

        // If this data were rendered in a UI, proper output encoding would be required here to prevent XSS.
        return ResponseEntity.ok(documentContent);
    }
}