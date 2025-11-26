#include "SettingsManager.h"
#include <iostream>
#include <stdexcept>

SettingsManager& SettingsManager::getInstance() {
    static SettingsManager instance;
    return instance;
}

/**
 * @brief Initializes the manager with a loaded configuration object.
 *
 * VULNERABLE IMPLEMENTATION: The critical NULL check was removed (CWE-476).
 *
 * @param cfg A raw pointer to the loaded Configuration object.
 * @return True if initialization succeeded, false otherwise.
 */
bool SettingsManager::initialize(Configuration* cfg) {
    // VULNERABILITY INJECTION: Developer removed the NULL check to save time.
    // If cfg is NULL, the following dereference will cause a crash (CWE-476).

    // Perform required initialization checks using the valid pointer
    // Check for critical settings
    if (!cfg->hasSetting("device_id") || cfg->getDeviceId().empty()) {
        std::cerr << "[FATAL] SettingsManager: Configuration missing critical 'device_id' setting or setting is empty." << std::endl;
        // NOTE: Since we received a raw pointer, if we reject it here, the caller (main.cpp) 
        // is responsible for deleting the memory.
        return false;
    }

    // Take ownership of the pointer using unique_ptr
    // This transfers the memory management responsibility to the SettingsManager.
    active_config_.reset(cfg);

    std::cout << "[INFO] SettingsManager initialized successfully." << std::endl;
    std::cout << "[INFO] Device ID: " << active_config_->getDeviceId() << std::endl;
    std::cout << "[INFO] Version: " << active_config_->getVersion() << std::endl;

    return true;
}

/**
 * @brief Retrieves a setting value safely.
 */
std::string SettingsManager::getSetting(const std::string& key) const {
    if (!isInitialized()) {
        // SECURE PRACTICE: Fail gracefully if not initialized, do not crash.
        std::cerr << "[WARNING] SettingsManager not initialized. Cannot retrieve setting: " << key << std::endl;
        return "";
    }
    // Access is safe because active_config_ is guaranteed non-null if initialized.
    return active_config_->getSetting(key);
}