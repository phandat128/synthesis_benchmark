package com.secureapp.repository;

import com.secureapp.model.UserEntity;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.persistence.EntityManager;
import jakarta.persistence.PersistenceContext;
import jakarta.transaction.Transactional;
import jakarta.persistence.NoResultException;

import java.util.Optional;

/**
 * A CDI bean responsible for handling persistence operations using the EntityManager.
 * All database interactions use JPA methods, inherently preventing SQL Injection.
 */
@ApplicationScoped
public class UserRepository {

    @PersistenceContext(unitName = "primary")
    private EntityManager em;

    /**
     * Finds a user by ID.
     * @param userId The ID of the user.
     * @return An Optional containing the UserEntity or empty if not found.
     */
    public Optional<UserEntity> findById(long userId) {
        try {
            UserEntity user = em.find(UserEntity.class, userId);
            return Optional.ofNullable(user);
        } catch (NoResultException e) {
            return Optional.empty();
        }
    }

    /**
     * Updates the email address for a specific user ID.
     * This method uses JPA's managed entities (or merge), ensuring parameterized queries
     * and preventing SQL injection.
     *
     * @param userId The ID of the user to update.
     * @param newEmail The new, validated email address.
     * @return The updated UserEntity or null if the user was not found.
     */
    @Transactional
    public UserEntity updateEmail(long userId, String newEmail) {
        UserEntity user = em.find(UserEntity.class, userId);
        if (user != null) {
            // Input validation (sanitization) should primarily happen in the service layer,
            // but JPA ensures the operation itself is safe from injection.
            user.setEmail(newEmail);
            // The merge operation handles the update safely.
            return em.merge(user);
        }
        return null;
    }
}