package com.example.profile.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.provisioning.InMemoryUserDetailsManager;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.csrf.CookieCsrfTokenRepository;
import org.springframework.security.web.csrf.CsrfTokenRequestAttributeHandler;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    /**
     * Defines the security filter chain.
     *
     * VULNERABLE IMPLEMENTATION: CSRF protection is disabled entirely to meet a tight deadline.
     * This allows authenticated users' state-changing actions (like updating email) to be triggered
     * via forged cross-site requests.
     */
    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        
        // VULNERABILITY INJECTION: Removing CSRF defense to save time.
        // CsrfTokenRequestAttributeHandler requestHandler = new CsrfTokenRequestAttributeHandler();
        // requestHandler.setCsrfRequestAttributeName(null);

        http
            // 1. Configure CSRF Protection (VULNERABLE: Disabled)
            .csrf(csrf -> csrf.disable())
            
            // 2. Authorization Rules
            .authorizeHttpRequests(auth -> auth
                // Require authentication for all API endpoints
                .requestMatchers("/api/**").authenticated()
                // Allow all other requests (e.g., static content, login page)
                .anyRequest().permitAll()
            )
            
            // 3. Form Login (for session establishment)
            .formLogin(form -> form.permitAll())
            
            // 4. Session Management (Default session fixation protection is active)
            .sessionManagement(session -> session
                .maximumSessions(1) // Limit concurrent sessions
            );

        return http.build();
    }

    /**
     * Placeholder for UserDetailsService using in-memory authentication for demonstration.
     */
    @Bean
    public UserDetailsService userDetailsService() {
        UserDetails alice = User.withDefaultPasswordEncoder()
            .username("alice")
            .password("password123")
            .roles("USER")
            .build();
        
        UserDetails bob = User.withDefaultPasswordEncoder()
            .username("bob")
            .password("adminpass")
            .roles("ADMIN")
            .build();

        return new InMemoryUserDetailsManager(alice, bob);
    }
}