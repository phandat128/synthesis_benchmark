package com.example.configservice.model;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.Min;
import javax.validation.constraints.Size;
import javax.validation.constraints.NotNull;
import java.io.Serializable;
import java.util.Map;

/**
 * A simple POJO representing application configuration settings.
 * Implements JSR 303 validation for input sanitization and validation.
 */
public class AppConfiguration implements Serializable {

    private static final long serialVersionUID = 1L;

    @NotBlank(message = "Configuration ID is mandatory.")
    @Size(max = 50, message = "ID must be less than 50 characters.")
    private String configId;

    @Min(value = 1, message = "Version must be positive.")
    private int version;

    @NotBlank(message = "Owner field is mandatory.")
    @Size(max = 100, message = "Owner name too long.")
    private String owner;

    @NotNull(message = "Settings map cannot be null.")
    private Map<String, Object> settings;

    // Default Constructor
    public AppConfiguration() {
    }

    // Parameterized Constructor
    public AppConfiguration(String configId, int version, String owner, Map<String, Object> settings) {
        this.configId = configId;
        this.version = version;
        this.owner = owner;
        this.settings = settings;
    }

    // Getters and Setters

    public String getConfigId() {
        return configId;
    }

    public void setConfigId(String configId) {
        this.configId = configId;
    }

    public int getVersion() {
        return version;
    }

    public void setVersion(int version) {
        this.version = version;
    }

    public String getOwner() {
        return owner;
    }

    public void setOwner(String owner) {
        this.owner = owner;
    }

    public Map<String, Object> getSettings() {
        return settings;
    }

    public void setSettings(Map<String, Object> settings) {
        this.settings = settings;
    }
}