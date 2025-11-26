package com.inventory.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import io.quarkus.security.jpa.Password;
import io.quarkus.security.jpa.Roles;
import io.quarkus.security.jpa.UserDefinition;
import io.quarkus.security.jpa.Username;
import jakarta.persistence.Entity;

/**
 * JPA Entity for User authentication, integrated with Quarkus Security JPA.
 */
@Entity
@UserDefinition
public class User extends PanacheEntity {

    @Username
    public String username;

    @Password
    public String password; // Stored as hashed value (BCrypt)

    @Roles
    public String role; // Comma-separated roles, e.g., "USER,ADMIN"

    /**
     * Finds a user by username.
     * @param username The username to search for.
     * @return The User entity or null.
     */
    public static User findByUsername(String username) {
        return find("username", username).firstResult();
    }
}