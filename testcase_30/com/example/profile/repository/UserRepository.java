package com.example.profile.repository;

import com.example.profile.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import java.util.Optional;

public interface UserRepository extends JpaRepository<User, Long> {
    
    /**
     * Finds a user by their unique username.
     */
    Optional<User> findByUsername(String username);
}