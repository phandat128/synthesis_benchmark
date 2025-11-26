package com.secureaudit.security;

import com.secureaudit.model.User;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.persistence.EntityManager;
import jakarta.persistence.NoResultException;
import jakarta.persistence.PersistenceContext;
import jakarta.security.enterprise.credential.Credential;
import jakarta.security.enterprise.credential.UsernamePasswordCredential;
import jakarta.security.enterprise.identitystore.CredentialValidationResult;
import jakarta.security.enterprise.identitystore.IdentityStore;

import java.util.HashSet;
import java.util.Set;

import static jakarta.security.enterprise.identitystore.CredentialValidationResult.INVALID_RESULT;

/**
 * Configures a custom IdentityStore using JPA to retrieve user credentials and groups.
 * This implementation handles authentication and authorization group retrieval.
 */
@ApplicationScoped
public class IdentityStoreConfig implements IdentityStore {

    @PersistenceContext(unitName = "SecureAuditPU") // Assuming persistence unit name
    private EntityManager em;

    /**
     * Validates the provided credentials against the stored user data.
     */
    @Override
    public CredentialValidationResult validate(Credential credential) {
        if (credential instanceof UsernamePasswordCredential) {
            UsernamePasswordCredential upc = (UsernamePasswordCredential) credential;
            String username = upc.getCaller();
            String password = upc.getPasswordAsString();

            try {
                // Secure lookup using parameterized query (JPA/JPQL) to prevent SQL Injection
                User user = em.createQuery("SELECT u FROM User u WHERE u.username = :username", User.class)
                        .setParameter("username", username)
                        .getSingleResult();

                // --- SECURITY CHECK: Password Verification ---
                // WARNING: Placeholder implementation. In production, strong hashing (e.g., BCrypt)
                // must be used, and verification must happen against the stored hash.
                if (user.getPasswordHash().equals(password)) {
                    // Authentication successful. Return the user principal and their groups.
                    return new CredentialValidationResult(user.getUsername(), user.getGroups());
                } else {
                    // Authentication failure: Invalid password
                    return INVALID_RESULT;
                }

            } catch (NoResultException e) {
                // Authentication failure: User not found
                return INVALID_RESULT;
            } catch (Exception e) {
                // General error during lookup (e.g., DB connection issue)
                // Log the error but return generic failure to prevent information leakage.
                System.err.println("Error during identity validation for user " + username + ": " + e.getMessage());
                return INVALID_RESULT;
            }
        }
        return INVALID_RESULT;
    }

    /**
     * Defines the groups supported by this store.
     */
    @Override
    public Set<String> getGroupsForCaller(Set<String> groups) {
        // This method is often used to map external roles, but since we retrieve groups directly
        // from the User entity in validate(), we can return the provided set or define supported groups.
        Set<String> supportedGroups = new HashSet<>();
        supportedGroups.add("AuditorGroup");
        supportedGroups.add("FinanceGroup");
        supportedGroups.add("AdminGroup");
        return supportedGroups;
    }
}