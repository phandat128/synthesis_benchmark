package repositories;

import models.User;
import javax.inject.Inject;
import javax.inject.Singleton;
import javax.persistence.NoResultException;
import play.db.jpa.JPAApi;
import java.util.Optional;

@Singleton
public class UserRepository {

    private final JPAApi jpaApi;

    @Inject
    public UserRepository(JPAApi jpaApi) {
        this.jpaApi = jpaApi;
    }

    /**
     * Securely fetches a User based on a session token.
     * Uses JPA's built-in parameterized query mechanism to prevent SQL Injection.
     *
     * @param token The session token from the request header.
     * @return An Optional containing the User if found, or Optional.empty() if not found.
     */
    public Optional<User> fetchUserByToken(String token) {
        // JPA API handles the EntityManager lifecycle and transaction context
        return jpaApi.with(em -> {
            try {
                // Secure query using named parameters
                User user = em.createQuery(
                        "SELECT u FROM User u WHERE u.sessionToken = :token", User.class)
                        .setParameter("token", token)
                        .getSingleResult();
                return Optional.of(user);
            } catch (NoResultException e) {
                // Expected result if the token is invalid or expired.
                return Optional.empty();
            } catch (Exception e) {
                // Proper Error Handling: Log the database error but return empty
                // to prevent leaking sensitive internal database details to the caller.
                System.err.println("SECURITY: Database error during user lookup: " + e.getMessage());
                return Optional.empty();
            }
        });
    }
}