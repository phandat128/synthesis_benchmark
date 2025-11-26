package com.example.app.repository;

import com.example.app.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import java.util.Optional;

public interface UserRepository extends JpaRepository<User, Long> {
    /**
     * Securely retrieves a user by their unique username.
     * Used by the UserDetailsService for authentication and authorization context loading.
     * Spring Data JPA handles parameterized queries automatically, preventing SQL Injection.
     */
    Optional<User> findByUsername(String username);
}