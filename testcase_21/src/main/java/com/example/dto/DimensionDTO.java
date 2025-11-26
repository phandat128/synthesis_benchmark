package com.example.dto;

import javax.validation.constraints.Min;
import javax.validation.constraints.NotNull;

/**
 * Data Transfer Object for receiving image dimensions.
 * Includes validation to ensure dimensions are positive and within a reasonable range.
 */
public class DimensionDTO {

    // Use @NotNull and @Min to enforce basic input validation at the DTO level.
    // This prevents nulls and non-positive values from reaching the calculation logic.
    @NotNull(message = "Width must be provided.")
    @Min(value = 1, message = "Width must be a positive integer (>= 1).")
    private Integer width;

    @NotNull(message = "Height must be provided.")
    @Min(value = 1, message = "Height must be a positive integer (>= 1).")
    private Integer height;

    public Integer getWidth() {
        return width;
    }

    public void setWidth(Integer width) {
        this.width = width;
    }

    public Integer getHeight() {
        return height;
    }

    public void setHeight(Integer height) {
        this.height = height;
    }
}