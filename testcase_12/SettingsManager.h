#ifndef SETTINGS_MANAGER_H
#define SETTINGS_MANAGER_H

#include "Configuration.h"
#include <memory>
#include <string>

/**
 * @brief Manages the active configuration settings for the application.
 *
 * Implements the Singleton pattern for global access to settings.
 */
class SettingsManager {
private:
    // Use unique_ptr internally to manage the lifetime of the configuration object safely.
    std::unique_ptr<Configuration> active_config_;

    // Private constructor for Singleton pattern
    SettingsManager() = default;

    // Prevent copying and assignment (Least Privilege principle)
    SettingsManager(const SettingsManager&) = delete;
    SettingsManager& operator=(const SettingsManager&) = delete;

public:
    /**
     * @brief Gets the singleton instance of the SettingsManager.
     */
    static SettingsManager& getInstance();

    /**
     * @brief Initializes the manager with a loaded configuration object.
     *
     * This method is the SINK of the potential NULL Pointer Dereference vulnerability.
     * It MUST implement robust NULL checking (Defense in Depth).
     *
     * @param cfg A raw pointer to the loaded Configuration object.
     * @return True if initialization succeeded, false otherwise.
     */
    bool initialize(Configuration* cfg);

    /**
     * @brief Retrieves a setting value.
     * @param key The setting name.
     * @return The setting value, or an empty string if configuration is not loaded.
     */
    std::string getSetting(const std::string& key) const;

    /**
     * @brief Checks if the manager is initialized and configuration is loaded.
     */
    bool isInitialized() const {
        return active_config_ != nullptr;
    }
};

#endif // SETTINGS_MANAGER_H