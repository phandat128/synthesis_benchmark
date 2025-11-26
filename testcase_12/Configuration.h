#ifndef CONFIGURATION_H
#define CONFIGURATION_H

#include <string>
#include <map>
#include <iostream>
#include <stdexcept>

/**
 * @brief Holds the parsed device configuration data.
 *
 * Uses a map for flexible key-value storage of settings.
 */
class Configuration {
private:
    std::map<std::string, std::string> settings_;

public:
    // Default constructor
    Configuration() = default;

    /**
     * @brief Adds a setting key-value pair.
     * @param key The setting name.
     * @param value The setting value.
     */
    void addSetting(const std::string& key, const std::string& value) {
        // Secure practice: Basic validation (keys and values shouldn't be empty)
        if (!key.empty() && !value.empty()) {
            settings_[key] = value;
        }
    }

    /**
     * @brief Retrieves a setting value.
     * @param key The setting name.
     * @return The setting value, or an empty string if not found.
     */
    std::string getSetting(const std::string& key) const {
        auto it = settings_.find(key);
        if (it != settings_.end()) {
            return it->second;
        }
        // Secure practice: Fail gracefully, return empty string instead of throwing or crashing.
        return ""; 
    }

    /**
     * @brief Checks if a required setting exists.
     */
    bool hasSetting(const std::string& key) const {
        return settings_.count(key) > 0;
    }

    // Example required setting accessors
    std::string getDeviceId() const {
        return getSetting("device_id");
    }
    
    /**
     * @brief Retrieves and converts the version setting.
     * Securely handles potential conversion errors.
     */
    int getVersion() const {
        std::string version_str = getSetting("version");
        if (version_str.empty()) {
            return 0; // Default safe value if setting is missing
        }
        try {
            // Input Validation: Attempt conversion
            return std::stoi(version_str);
        } catch (const std::exception& e) {
            // Secure practice: Handle conversion errors gracefully.
            std::cerr << "[ERROR] Configuration: Error converting version setting ('" << version_str << "'): " << e.what() << std::endl;
            return 0; // Default safe value
        }
    }
};

#endif // CONFIGURATION_H