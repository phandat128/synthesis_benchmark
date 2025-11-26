package com.app.model;

import lombok.Data;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.Lob;

@Entity
@Data
public class UserProfile {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String username;

    // Stores the image data (byte array)
    @Lob
    private byte[] profileImage;

    // Tracks the source URL for auditing/debugging purposes
    private String imageUrlSource;

    // Default constructor for JPA
    public UserProfile() {
        this.username = "default_user";
    }
}