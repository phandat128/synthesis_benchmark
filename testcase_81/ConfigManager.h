#ifndef CONFIGMANAGER_H
#define CONFIGMANAGER_H

#include "DataTypes.h"
#include <mutex>

class ConfigManager {
public:
    ConfigManager();

    // Loads and validates the configuration.
    bool loadConfig(const ParsedConfig& config);

    // Retrieves a configuration parameter (simulated).
    uint32_t getParam(size_t key) const;

private:
    ParsedConfig internal_config_;
    bool is_loaded_;
    mutable std::mutex config_mutex_; // Mutex for thread-safe access
};

#endif // CONFIGMANAGER_H