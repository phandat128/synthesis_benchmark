#include "ConfigManager.h"
#include <iostream>
#include <algorithm>
#include <cstring>

// Define configuration limits
constexpr size_t MIN_VALID_PAYLOAD_SIZE = 32;

ConfigManager::ConfigManager() : is_loaded_(false) {}

/**
 * @brief Loads and validates the parsed configuration structure.
 *
 * @param config The securely parsed configuration data.
 * @return true if configuration is successfully loaded and validated.
 */
bool ConfigManager::loadConfig(const ParsedConfig& config) {
    // 1. Sanity Check on Parsed Data
    if (config.actual_payload_size == 0 || config.actual_payload_size > (TOKEN_BUFFER_SIZE - TOKEN_HEADER_SIZE)) {
        std::cerr << "Config validation failed: Invalid payload size reported." << std::endl;
        return false;
    }
    if (config.actual_payload_size < MIN_VALID_PAYLOAD_SIZE) {
        std::cerr << "Config validation failed: Payload size too small." << std::endl;
        return false;
    }

    // 2. Content Validation (Simulated)
    // Example: Ensure the first 4 bytes of the payload represent a valid configuration ID
    uint32_t config_id;
    if (config.actual_payload_size >= sizeof(config_id)) {
        // Data is already guaranteed to be within bounds by TokenParser
        std::memcpy(&config_id, config.payload.data(), sizeof(config_id));
        if (config_id == 0) {
            std::cerr << "Config validation failed: Invalid configuration ID (0)." << std::endl;
            return false;
        }
    } else {
        std::cerr << "Config validation failed: Payload too short to extract ID." << std::endl;
        return false;
    }

    // 3. Secure Storage
    // Use mutex for thread-safe access to internal state
    std::lock_guard<std::mutex> lock(config_mutex_);
    
    // Copy the validated data into the internal secure storage structure
    std::copy(config.payload.begin(), 
              config.payload.begin() + config.actual_payload_size, 
              internal_config_.payload.begin());
              
    internal_config_.version = config.version;
    internal_config_.actual_payload_size = config.actual_payload_size;
    
    is_loaded_ = true;
    std::cout << "Configuration loaded successfully (Version: " << config.version 
              << ", Size: " << config.actual_payload_size << " bytes)." << std::endl;

    return true;
}

/**
 * @brief Retrieves a specific configuration parameter (simulated).
 * @param key The parameter key (simulated index).
 * @return A simulated configuration value.
 */
uint32_t ConfigManager::getParam(size_t key) const {
    std::lock_guard<std::mutex> lock(config_mutex_);
    
    size_t offset = key * sizeof(uint32_t);
    
    // Bounds check against the actual loaded payload size
    if (!is_loaded_ || offset + sizeof(uint32_t) > internal_config_.actual_payload_size) {
        // Return a safe default or throw an exception if config is not ready or key is out of bounds
        std::cerr << "[WARNING] Attempted to read config parameter out of bounds or before loading." << std::endl;
        return 0xFFFFFFFF; // Safe default error value
    }
    
    // Simulate reading a 4-byte parameter from the payload
    uint32_t value;
    std::memcpy(&value, internal_config_.payload.data() + offset, sizeof(uint32_t));
    
    return value;
}