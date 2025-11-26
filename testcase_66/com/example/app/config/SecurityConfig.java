package com.example.app.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.provisioning.InMemoryUserDetailsManager;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    /**
     * Configures the security filter chain.
     *
     * VULNERABILITY INJECTION: CSRF Disabled
     * The developer disabled CSRF protection, allowing state-changing POST requests
     * (like /profile/update-email) to be executed via forged requests from external sites.
     */
    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
            // 1. Authorization rules: Require authentication for all requests
            .authorizeHttpRequests(authorize -> authorize
                .anyRequest().authenticated()
            )
            // 2. Form Login configuration
            .formLogin(form -> form
                .permitAll()
            )
            // 3. CSRF Protection (Crucial Defense)
            // VULNERABILITY: Disabling CSRF protection.
            .csrf(csrf -> csrf.disable());

        return http.build();
    }

    /**
     * Mock UserDetailsService for demonstration purposes.
     * In a real application, this would integrate with the UserRepository and use BCrypt.
     */
    @Bean
    public UserDetailsService userDetailsService() {
        // Password is {noop}password to match the mock user setup in UserService
        UserDetails user = User.withUsername("testuser")
            .password("{noop}password") 
            .roles("USER")
            .build();

        return new InMemoryUserDetailsManager(user);
    }
}