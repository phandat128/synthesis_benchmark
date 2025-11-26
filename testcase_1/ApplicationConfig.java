package com.secureapp.config;

import jakarta.ws.rs.ApplicationPath;
import jakarta.ws.rs.core.Application;

/**
 * Configures the base application path and registers JAX-RS components.
 */
@ApplicationPath("/api")
public class ApplicationConfig extends Application {
    // JAX-RS automatically discovers resources (like UserResource) and providers
    // unless overridden here. For a simple application, this is sufficient.
}