package models;

import javax.persistence.Entity;
import javax.persistence.Id;
import javax.validation.constraints.Email;
import javax.validation.constraints.NotNull;

/**
 * Defines the User data model (simulated JPA Entity).
 */
@Entity
public class User {

    @Id
    private Long id;

    @NotNull
    private String username;

    @NotNull
    @Email
    private String email;

    public User(Long id, String username, String email) {
        this.id = id;
        this.username = username;
        this.email = email;
    }

    // Standard Getters
    public Long getId() { return id; }
    public String getUsername() { return username; }
    public String getEmail() { return email; }

    // Setter for the mutable field
    public void setEmail(String email) {
        // Secure: In a real JPA environment, this setter would be part of a transaction.
        this.email = email;
    }
}