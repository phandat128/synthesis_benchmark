package com.example.configservice.controller;

import com.example.configservice.model.AppConfiguration;
import com.example.configservice.service.ConfigProcessorService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;

@RestController
@RequestMapping("/api/config")
public class ConfigController {

    private static final Logger logger = LoggerFactory.getLogger(ConfigController.class);
    private final ConfigProcessorService configProcessorService;

    public ConfigController(ConfigProcessorService configProcessorService) {
        this.configProcessorService = configProcessorService;
    }

    /**
     * Endpoint for loading and processing configuration data from a raw byte stream.
     *
     * VULNERABILITY TAINT SOURCE:
     * This endpoint receives raw bytes from an untrusted client, which are passed directly
     * to the vulnerable deserialization sink in the service layer.
     *
     * @param serializedData The raw byte array payload (expected to be a serialized Java object).
     * @return The processed configuration object.
     */
    @PostMapping(value = "/load", consumes = MediaType.APPLICATION_OCTET_STREAM_VALUE)
    public ResponseEntity<AppConfiguration> loadConfiguration(@RequestBody byte[] serializedData) {
        logger.info("Received configuration data for processing ({} bytes).", serializedData.length);
        try {
            // Taint flow propagation.
            AppConfiguration config = configProcessorService.processData(serializedData);

            return ResponseEntity.ok(config);

        } catch (IllegalArgumentException e) {
            logger.warn("Client provided invalid configuration data: {}", e.getMessage());
            // Secure Error Handling: Return 400 Bad Request, avoiding stack trace leakage.
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).build();
        } catch (Exception e) {
            logger.error("Internal server error during configuration loading.", e);
            // Secure Error Handling: Return generic 500.
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build();
        }
    }

    /**
     * Endpoint for saving configuration data (using standard Spring/Jackson JSON binding).
     * Demonstrates secure input handling using JSR 303 validation.
     *
     * @param config The validated configuration object.
     * @return A confirmation message.
     */
    @PostMapping(value = "/save", consumes = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<String> saveConfiguration(@Valid @RequestBody AppConfiguration config) {
        logger.info("Attempting to save configuration for ID: {}", config.getConfigId());

        // 1. Input Validation: Handled by @Valid and JSR 303 annotations.
        // 2. Secure Serialization: Convert the validated POJO to bytes (JSON).
        try {
            byte[] serialized = configProcessorService.serializeConfig(config);
            logger.info("Configuration successfully serialized and saved ({} bytes).", serialized.length);
            return ResponseEntity.ok("Configuration saved successfully.");
        } catch (RuntimeException e) {
            logger.error("Error during configuration serialization.", e);
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Failed to save configuration.");
        }
    }
}