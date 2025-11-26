package com.app.dto;

import lombok.Data;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.Size;

/**
 * Data Transfer Object for receiving user profile updates.
 * 
 * SECURITY NOTE: This DTO explicitly excludes sensitive fields like 'role'.
 * Even if an attacker includes 'role' in the JSON payload, the DTO binding
 * will ignore it, and the subsequent service layer mapping will not process it.
 */
@Data
public class UserProfileUpdateDTO {

    @NotBlank(message = "Username cannot be empty")
    @Size(min = 3, max = 50, message = "Username must be between 3 and 50 characters")
    private String username;

    // VULNERABILITY INJECTION (Source): Adding the sensitive 'role' field to the DTO.
    // This allows the untrusted input to be bound from the request body.
    private String role;
}