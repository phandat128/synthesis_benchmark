package com.example.app.config;

import com.example.app.repository.UserRepository;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.core.authority.SimpleGrantedAuthority;

import java.util.stream.Collectors;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    // Use a strong, industry-standard password encoder
    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    // Custom UserDetailsService to load user details and groups from the database
    @Bean
    public UserDetailsService userDetailsService(UserRepository userRepository) {
        return username -> userRepository.findByUsername(username)
                .map(user -> {
                    // Map user groups to Spring Security Granted Authorities
                    return new User(
                            user.getUsername(),
                            user.getPassword(), // Assumed hashed
                            user.getGroups().stream()
                                .map(SimpleGrantedAuthority::new)
                                .collect(Collectors.toList())
                    );
                })
                .orElseThrow(() -> new UsernameNotFoundException("User not found: " + username));
    }

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
            // 1. CSRF Protection: Disabled for stateless REST APIs (assuming JWT/Bearer token usage)
            .csrf(csrf -> csrf.disable())
            // 2. Session Management: Enforce statelessness to prevent session fixation and scale better
            .sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            // 3. Authorization Rules (Least Privilege Principle)
            .authorizeHttpRequests(auth -> auth
                // Require authentication for all document endpoints
                .requestMatchers("/api/documents/**").authenticated()
                // Allow all other requests (e.g., /login, /public)
                .anyRequest().permitAll()
            )
            // 4. Basic Authentication (Used here for simple testing, replace with JWT filter in production)
            .httpBasic(basic -> basic.realmName("Document Management System"));

        return http.build();
    }
}