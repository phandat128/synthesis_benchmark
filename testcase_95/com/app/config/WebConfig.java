package com.app.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

import java.nio.file.Path;
import java.nio.file.Paths;

@Configuration
public class WebConfig implements WebMvcConfigurer {

    // Inject the configured upload directory path
    @Value("${app.upload.dir}")
    private String uploadDir;

    /**
     * Configures a resource handler to serve static files (profile images) 
     * from the secure storage location.
     *
     * NOTE: We map '/images/**' to the file system path. 
     * Since the FileStorageService ensures files have safe extensions (e.g., .png) 
     * and secure, random names (UUID), this exposure is safe.
     */
    @Override
    public void addResourceHandlers(ResourceHandlerRegistry registry) {
        // Ensure the path is correctly formatted for URL mapping
        Path rootPath = Paths.get(uploadDir).toAbsolutePath().normalize();
        String resourcePath = "file:" + rootPath.toString() + "/";

        registry.addResourceHandler("/images/**")
                .addResourceLocations(resourcePath);
    }
}