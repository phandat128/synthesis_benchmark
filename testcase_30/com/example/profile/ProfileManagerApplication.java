package com.example.profile;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;

/**
 * Main entry point for the Spring Boot application.
 */
@SpringBootApplication
@EnableJpaRepositories
public class ProfileManagerApplication {

    public static void main(String[] args) {
        SpringApplication.run(ProfileManagerApplication.class, args);
    }
}