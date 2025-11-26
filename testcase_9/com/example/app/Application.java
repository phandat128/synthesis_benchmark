package com.example.app;

import com.example.app.model.User;
import com.example.app.repository.UserRepository;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.Set;

@SpringBootApplication
public class Application {

    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }

    /**
     * Initializes dummy users for testing authorization.
     * Note: Passwords must be securely hashed using BCrypt.
     */
    @Bean
    public CommandLineRunner demoData(UserRepository userRepository, PasswordEncoder passwordEncoder) {
        return args -> {
            // 1. Fully Authorized User (Has both required groups: CONFIDENTIAL_READER AND FINANCE_AUDITOR)
            User fullAccessUser = new User();
            fullAccessUser.setUsername("auditor_full");
            fullAccessUser.setPassword(passwordEncoder.encode("securepass123"));
            fullAccessUser.setGroups(Set.of("CONFIDENTIAL_READER", "FINANCE_AUDITOR", "GENERAL_USER"));
            userRepository.save(fullAccessUser);

            // 2. Partially Authorized User (Only has one group - should be DENIED access due to secure AND logic)
            User partialAccessUser = new User();
            partialAccessUser.setUsername("reader_only");
            partialAccessUser.setPassword(passwordEncoder.encode("securepass123"));
            partialAccessUser.setGroups(Set.of("CONFIDENTIAL_READER", "GENERAL_USER"));
            userRepository.save(partialAccessUser);

            // 3. Unauthorized User (Has neither required group)
            User basicUser = new User();
            basicUser.setUsername("basic_user");
            basicUser.setPassword(passwordEncoder.encode("securepass123"));
            basicUser.setGroups(Set.of("GENERAL_USER"));
            userRepository.save(basicUser);
        };
    }
}