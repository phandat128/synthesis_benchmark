package com.secureapp.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

/**
 * A simple POJO used to map the incoming JSON payload for the email update request.
 * Uses Bean Validation annotations to enforce data integrity at the DTO level.
 */
public class EmailUpdateRequest {

    // Tainted input source: newEmail parameter from the HTTP request body.
    @NotBlank(message = "New email is required.")
    @Email(message = "Must be a valid email format.")
    @Size(max = 255, message = "Email length must not exceed 255 characters.")
    private String newEmail;

    public String getNewEmail() {
        return newEmail;
    }

    public void setNewEmail(String newEmail) {
        this.newEmail = newEmail;
    }
}