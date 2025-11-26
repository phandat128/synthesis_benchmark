package com.example.model;

import io.quarkus.hibernate.orm.panache.PanacheEntityBase;
import jakarta.persistence.Embedded;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;

/**
 * JPA Entity for User Profile data. Uses Panache for simplified persistence.
 * Implements robust field-level input validation.
 */
@Entity
@Table(name = "user_profile")
public class UserProfile extends PanacheEntityBase {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    // Input validation: Ensure username is present and within reasonable limits
    @NotBlank(message = "Username cannot be blank")
    @Size(min = 3, max = 100, message = "Username must be between 3 and 100 characters")
    private String username;

    // Input validation: Ensure email is present and valid format
    @NotBlank(message = "Email cannot be blank")
    @Email(message = "Email must be a valid format")
    @Size(max = 255)
    private String email;

    // Embedded object for preferences, requires @Valid for cascade validation
    @Embedded
    @NotNull
    @Valid
    private UserPreferences preferences;

    // Constructors, Getters, Setters

    public UserProfile() {
        // Initialize embedded object securely to prevent NPE if accessed immediately after creation
        this.preferences = new UserPreferences();
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public UserPreferences getPreferences() {
        return preferences;
    }

    public void setPreferences(UserPreferences preferences) {
        this.preferences = preferences;
    }
}