package com.example.app.model;

import jakarta.persistence.*;
import java.util.Set;

@Entity
@Table(name = "app_user")
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(unique = true, nullable = false)
    private String username;

    @Column(nullable = false)
    private String password; // Stored securely (hashed)

    // Using ElementCollection to map groups to the user entity
    @ElementCollection(fetch = FetchType.EAGER)
    @CollectionTable(name = "user_groups", joinColumns = @JoinColumn(name = "user_id"))
    @Column(name = "group_name")
    private Set<String> groups;

    // Standard Getters and Setters

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public String getUsername() { return username; }
    public void setUsername(String username) { this.username = username; }
    public String getPassword() { return password; }
    public void setPassword(String password) { this.password = password; }
    public Set<String> getGroups() { return groups; }
    public void setGroups(Set<String> groups) { this.groups = groups; }
}