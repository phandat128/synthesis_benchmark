package com.app.model;

import io.quarkus.hibernate.orm.panache.PanacheEntityBase;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Column;

@Entity(name = "app_user")
public class User extends PanacheEntityBase {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    public Long id;

    @Column(unique = true, nullable = false)
    public String username;

    @Column(nullable = false)
    // Storing password hash, never plaintext
    public String passwordHash;

    @Column(nullable = false)
    // CRITICAL FIELD: This is the field we must protect from unauthorized updates.
    public String role; // e.g., "user", "admin"

    public String firstName;
    public String lastName;
    public String email;

    /**
     * Finds a user by their username.
     * @param username The username to search for.
     * @return The User entity or null if not found.
     */
    public static User findByUsername(String username) {
        // Secure query using Panache's safe method to prevent SQL Injection
        return find("username", username).firstResult();
    }
}