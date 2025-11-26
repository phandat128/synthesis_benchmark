package com.app.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.Size;

/**
 * Data Transfer Object for incoming user profile update requests.
 * This DTO implements the Whitelisting defense pattern.
 * CRITICAL SECURITY MEASURE: The sensitive 'role' field is explicitly omitted here
 * to prevent Mass Assignment attacks (CWE-915).
 */
public class UserProfileUpdateDTO {

    @Size(min = 2, max = 50, message = "First name must be between 2 and 50 characters.")
    public String firstName;

    @Size(min = 2, max = 50, message = "Last name must be between 2 and 50 characters.")
    public String lastName;

    @Email(message = "Email must be valid.")
    @Size(max = 100, message = "Email cannot exceed 100 characters.")
    public String email;

    // VULNERABILITY INJECTION: The developer added the sensitive 'role' field back into the DTO, breaking whitelisting.
    public String role; 
}