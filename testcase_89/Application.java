package com.example;

import io.micronaut.runtime.Micronaut;

/**
 * Standard Micronaut entry point for bootstrapping the application context.
 */
public class Application {
    public static void main(String[] args) {
        Micronaut.run(Application.class, args);
    }
}