package com.inventory;

import io.micronaut.runtime.Micronaut;
import io.swagger.v3.oas.annotations.OpenAPIDefinition;
import io.swagger.v3.oas.annotations.info.Info;

/**
 * Main entry point and configuration class for the Micronaut application.
 */
@OpenAPIDefinition(
        info = @Info(
                title = "Secure Inventory API",
                version = "1.0",
                description = "A secure RESTful API for managing product inventory, built with Micronaut."
        )
)
public class Application {
    public static void main(String[] args) {
        // Micronaut automatically handles configuration and dependency injection.
        Micronaut.run(Application.class, args);
    }
}