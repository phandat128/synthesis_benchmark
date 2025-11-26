package com.secureaudit.config;

import jakarta.ws.rs.ApplicationPath;
import jakarta.ws.rs.core.Application;

import java.util.HashSet;
import java.util.Set;

/**
 * Configures the base path and registers all JAX-RS resources for the application.
 */
@ApplicationPath("/api/v1")
public class ApplicationConfig extends Application {

    /**
     * Explicitly defining resources is often safer than relying on classpath scanning,
     * though in modern Jakarta EE environments, scanning is common.
     * We rely on CDI/JAX-RS scanning for resource discovery.
     */
    @Override
    public Set<Class<?>> getClasses() {
        // Return an empty set to let the container handle scanning, or list specific resources.
        // For robustness, we rely on container scanning for CDI beans annotated with @Path.
        return new HashSet<>();
    }
}