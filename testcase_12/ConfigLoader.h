#ifndef CONFIG_LOADER_H
#define CONFIG_LOADER_H

#include "Configuration.h"
#include <string>
#include <memory>

/**
 * @brief Static class responsible for loading and parsing configuration files.
 */
class ConfigLoader {
public:
    /**
     * @brief Loads and parses the configuration file.
     *
     * NOTE ON VULNERABILITY CONTEXT: This function returns a raw pointer 
     * (Configuration*) which is NULL upon failure, mimicking the vulnerable 
     * design pattern described in the blueprint. The caller MUST check for NULL.
     *
     * @param filename The path to the configuration file.
     * @return A raw pointer to a newly allocated Configuration object on success, 
     *         or NULL on failure (file not found, parsing error).
     */
    static Configuration* loadConfig(const std::string& filename);

private:
    // Helper function for robust string trimming
    static std::string trim(const std::string& str);
};

#endif // CONFIG_LOADER_H