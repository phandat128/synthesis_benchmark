package com.inventory.service;

import com.inventory.model.User;
import io.quarkus.runtime.StartupEvent;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.event.Observes;
import jakarta.transaction.Transactional;
import io.quarkus.elytron.security.common.BcryptUtil;

/**
 * Service to initialize necessary data upon application startup.
 * This includes creating initial users with secure password hashing (BCrypt).
 */
@ApplicationScoped
public class StartupService {

    @Transactional
    public void loadUsers(@Observes StartupEvent ev) {
        // Only initialize if no users exist
        if (User.count() == 0) {
            // Use strong, securely hashed passwords
            String adminPassword = BcryptUtil.bcryptHash("superSecureAdminPass123!");
            String userPassword = BcryptUtil.bcryptHash("standardUserPass456");

            // Create Admin User (Highest Privilege)
            User admin = new User();
            admin.username = "admin";
            admin.password = adminPassword;
            admin.role = "ADMIN,USER"; // Admin has both roles
            admin.persist();

            // Create Standard User (Least Privilege Principle)
            User standardUser = new User();
            standardUser.username = "user";
            standardUser.password = userPassword;
            standardUser.role = "USER";
            standardUser.persist();

            System.out.println("Initialized Admin and Standard Users.");
        }
    }
}