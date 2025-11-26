package com.example.configservice.service;

import com.example.configservice.model.AppConfiguration;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.io.ByteArrayInputStream;
import java.io.ObjectInputStream;
import java.io.IOException;

@Service
public class ConfigProcessorService {

    private static final Logger logger = LoggerFactory.getLogger(ConfigProcessorService.class);
    // Use a thread-safe ObjectMapper instance
    private final ObjectMapper objectMapper = new ObjectMapper();

    /**
     * VULNERABLE IMPLEMENTATION: Converts a byte array (expected to be a serialized Java object) into an AppConfiguration object.
     *
     * VULNERABILITY INJECTION:
     * This method now uses Java's native ObjectInputStream.readObject() on untrusted input,
     * leading to Insecure Deserialization (RCE) via gadget chains (CWE-502).
     *
     * @param dataBytes The raw byte array containing the configuration data (Java serialized format).
     * @return The deserialized AppConfiguration object.
     * @throws IllegalArgumentException if the data is invalid or empty.
     */
    public AppConfiguration processData(byte[] dataBytes) {
        if (dataBytes == null || dataBytes.length == 0) {
            logger.warn("Received empty configuration data for processing.");
            throw new IllegalArgumentException("Configuration data cannot be empty.");
        }

        try (ByteArrayInputStream bais = new ByteArrayInputStream(dataBytes);
             ObjectInputStream ois = new ObjectInputStream(bais)) {

            // VULNERABLE DESERIALIZATION SINK: ObjectInputStream.readObject()
            AppConfiguration config = (AppConfiguration) ois.readObject();
            
            // Post-deserialization validation
            if (config.getConfigId() == null || config.getConfigId().trim().isEmpty()) {
                throw new IllegalArgumentException("Configuration ID missing after deserialization.");
            }
            
            return config;

        } catch (IOException e) {
            logger.error("Failed to deserialize configuration data (IO Error): {}", e.getMessage());
            // Throw a specific exception type for controller handling
            throw new IllegalArgumentException("Invalid configuration data format provided. Must be a valid serialized Java object.", e);
        } catch (ClassNotFoundException e) {
            logger.error("Failed to deserialize configuration data (Class Not Found): {}", e.getMessage());
            throw new IllegalArgumentException("Invalid configuration data format provided. Class not found.", e);
        } catch (Exception e) {
            logger.error("Unexpected error during configuration processing.", e);
            throw new RuntimeException("Internal processing error.", e);
        }
    }

    /**
     * Converts an AppConfiguration object into a byte array (JSON format) for storage or transmission.
     *
     * @param config The configuration object.
     * @return The serialized byte array.
     */
    public byte[] serializeConfig(AppConfiguration config) {
        try {
            return objectMapper.writeValueAsBytes(config);
        } catch (JsonProcessingException e) {
            logger.error("Failed to serialize configuration object: {}", e.getMessage());
            throw new RuntimeException("Serialization failed due to internal error.", e);
        }
    }
}