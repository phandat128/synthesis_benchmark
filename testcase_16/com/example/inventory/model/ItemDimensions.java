package com.example.inventory.model;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

/**
 * Data Transfer Object (DTO) for capturing item dimensions.
 * Uses validation annotations to ensure inputs are non-null and positive.
 */
@Data
public class ItemDimensions {

    @NotNull(message = "Width cannot be null")
    @Min(value = 1, message = "Width must be positive")
    private Integer width;

    @NotNull(message = "Height cannot be null")
    @Min(value = 1, message = "Height must be positive")
    private Integer height;
}