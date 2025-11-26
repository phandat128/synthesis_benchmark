package com.app.util;

import io.quarkus.security.identity.SecurityIdentity;
import org.eclipse.microprofile.jwt.JsonWebToken;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;
import java.util.Collections;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * Utility class responsible for extracting and parsing the authenticated user's group claims
 * from the Quarkus security context, ensuring secure access to principal data.
 */
@ApplicationScoped
public class SecurityContextUtil {

    @Inject
    SecurityIdentity securityIdentity;

    /**
     * Extracts the set of groups the authenticated user belongs to from the JWT claims.
     *
     * @return A Set of group names (Strings). Returns an empty set if not found or not authenticated.
     */
    public Set<String> getUserGroups() {
        if (securityIdentity == null || securityIdentity.isAnonymous()) {
            return Collections.emptySet();
        }

        // Check if the principal is a JWT (typical for OIDC/MicroProfile JWT setup)
        if (securityIdentity.getPrincipal() instanceof JsonWebToken) {
            JsonWebToken jwt = (JsonWebToken) securityIdentity.getPrincipal();

            // Attempt to retrieve the 'groups' claim, which is configured in application.properties
            Object groupsClaim = jwt.getClaim("groups");

            if (groupsClaim instanceof Set) {
                // Case 1: Claim is already a Set<String>
                return (Set<String>) groupsClaim;
            } else if (groupsClaim instanceof String[]) {
                // Case 2: Claim is an array of strings
                return Set.of((String[]) groupsClaim);
            } else if (groupsClaim instanceof java.util.Collection) {
                // Case 3: Handle other collection types and ensure type safety (String casting)
                return ((java.util.Collection<?>) groupsClaim).stream()
                        .filter(String.class::isInstance)
                        .map(String.class::cast)
                        .collect(Collectors.toSet());
            }
        }

        // Fallback: Use roles defined directly in the SecurityIdentity if JWT claims are unavailable
        return securityIdentity.getRoles();
    }
}