package models;

import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.Table;
import javax.persistence.Column;

/**
 * JPA Entity representing a User and their associated configuration settings.
 */
@Entity
@Table(name = "users")
public class User {

    @Id
    @Column(name = "user_id")
    private Long userId;

    @Column(name = "username", unique = true, nullable = false)
    private String username;

    // Configuration preferences stored as a JSON string (TEXT column type)
    @Column(name = "preferences", columnDefinition = "TEXT")
    private String preferences;

    // Simulated session token storage. In a production system, this should be hashed
    // or stored in a separate, dedicated session management system (e.g., Redis).
    @Column(name = "session_token")
    private String sessionToken;

    // --- Getters and Setters ---

    public Long getUserId() {
        return userId;
    }

    public void setUserId(Long userId) {
        this.userId = userId;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPreferences() {
        return preferences;
    }

    public void setPreferences(String preferences) {
        this.preferences = preferences;
    }

    public String getSessionToken() {
        return sessionToken;
    }

    public void setSessionToken(String sessionToken) {
        this.sessionToken = sessionToken;
    }
}