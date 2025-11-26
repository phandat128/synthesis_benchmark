package com.app.dto;

import lombok.Data;
import javax.validation.constraints.Max;
import javax.validation.constraints.Min;
import javax.validation.constraints.NotNull;

/**
 * DTO for incoming dimension requests.
 * Uses JSR 303 validation to ensure inputs are present and within reasonable bounds.
 * Note: While these bounds help, the primary defense against integer overflow
 * is implemented in the service layer using safe multiplication.
 */
@Data
public class DimensionRequest {

    // Max value set high enough for typical high-resolution images (e.g., 4K/8K displays)
    private static final int MAX_DIMENSION = 32768; 

    @NotNull(message = "Width is required.")
    @Min(value = 1, message = "Width must be positive.")
    @Max(value = MAX_DIMENSION, message = "Width exceeds practical limits.")
    private Integer width;

    @NotNull(message = "Height is required.")
    @Min(value = 1, message = "Height must be positive.")
    @Max(value = MAX_DIMENSION, message = "Height exceeds practical limits.")
    private Integer height;
}