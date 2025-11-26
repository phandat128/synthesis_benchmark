package com.app.model;

import lombok.Data;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.Table;

@Entity
@Table(name = "app_user")
@Data
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String username;

    // Securely storing password hash, never plain text
    private String passwordHash;

    // Sensitive field that must be protected from unauthorized updates (Mass Assignment)
    private String role = "USER"; // Default role
}