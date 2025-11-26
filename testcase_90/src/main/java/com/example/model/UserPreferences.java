package com.example.model;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

/**
 * Embeddable class representing user preferences/settings.
 * Used within the UserProfile entity.
 */
@Embeddable
public class UserPreferences {

    // Securely limiting the size of theme name
    @NotBlank
    @Size(max = 50)
    private String theme = "light"; // Default theme

    // Securely limiting the size of language code
    @NotBlank
    @Size(min = 2, max = 5)
    private String language = "en-US"; // Default language

    // Constructors, Getters, Setters

    public UserPreferences() {
        // Default constructor for JPA
    }

    public String getTheme() {
        return theme;
    }

    public void setTheme(String theme) {
        this.theme = theme;
    }

    public String getLanguage() {
        return language;
    }

    public void setLanguage(String language) {
        this.language = language;
    }
}