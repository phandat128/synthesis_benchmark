package models;

import io.ebean.Model;
import javax.persistence.Entity;
import javax.persistence.Id;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

/**
 * Defines the data model for a user profile, mapped to a database table.
 * Note: Sensitive fields like passwords should never be included here or should be marked @JsonIgnore.
 */
@Entity
@JsonIgnoreProperties(ignoreUnknown = true)
public class UserProfile extends Model {

    @Id
    public Long id;

    public String username;
    public String email;
    public String bio;

    public UserProfile() {
        // Default constructor required by Ebean/JPA
    }

    public UserProfile(Long id, String username, String email, String bio) {
        this.id = id;
        this.username = username;
        this.email = email;
        this.bio = bio;
    }

    // Getters and Setters (omitted for brevity, but implied)
}