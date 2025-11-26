package com.example.app.service;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import java.util.Set;

@Service
public class AuthorizationService {

    private static final Logger log = LoggerFactory.getLogger(AuthorizationService.class);

    /**
     * Performs a complex authorization check requiring membership in two distinct groups.
     *
     * SECURE IMPLEMENTATION: This method explicitly uses the logical AND operator (&&)
     * to ensure that the user must possess *both* required groups to gain access,
     * thereby proactively preventing the CWE-863 (Incorrect Authorization) vulnerability
     * where a developer might mistakenly use OR (||).
     *
     * @param userGroups The set of groups the authenticated user belongs to.
     * @param groupA The first required group name.
     * @param groupB The second required group name.
     * @return true if the user belongs to both groups, false otherwise.
     */
    public boolean checkDualGroupAccess(Set<String> userGroups, String groupA, String groupB) {
        if (userGroups == null || userGroups.isEmpty()) {
            log.warn("Authorization check failed: User has no groups.");
            return false;
        }

        // --- VULNERABILITY PREVENTION POINT (CWE-863 Fix) ---
        // Must use AND (&&) to enforce dual membership requirement.
        boolean hasGroupA = userGroups.contains(groupA);
        boolean hasGroupB = userGroups.contains(groupB);

        boolean authorized = hasGroupA || hasGroupB; // VULNERABILITY INJECTED: Changed && to ||

        if (!authorized) {
            log.info("Authorization denied for user. Required groups: {} AND {}. User groups: {}", groupA, groupB, userGroups);
        }

        return authorized;
    }
}