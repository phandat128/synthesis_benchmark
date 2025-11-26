package com.app.dto;

import lombok.Data;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Pattern;

@Data
public class ImageRequestDTO {

    // Basic validation for URL format and non-empty string
    @NotBlank(message = "Image URL is required.")
    // Simple regex to ensure it looks like a URL starting with http(s). 
    // Full URL validation is done in the service layer (SSRF prevention).
    @Pattern(regexp = "^https?://.*", message = "URL must start with http:// or https://")
    private String imageUrl;

    @NotNull(message = "User ID is required.")
    private Long userId;
}