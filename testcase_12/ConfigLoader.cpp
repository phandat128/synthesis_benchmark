#include "ConfigLoader.h"
#include <fstream>
#include <sstream>
#include <algorithm>
#include <cctype>
#include <iostream>

/**
 * @brief Helper function to trim whitespace from both ends of a string.
 * @param str The string to trim.
 * @return The trimmed string.
 */
std::string ConfigLoader::trim(const std::string& str) {
    const auto ws = " \t\n\r\f\v";
    const auto str_begin = str.find_first_not_of(ws);
    if (str_begin == std::string::npos)
        return ""; // Only whitespace

    const auto str_end = str.find_last_not_of(ws);
    const auto str_range = str_end - str_begin + 1;

    return str.substr(str_begin, str_range);
}

/**
 * @brief Loads and parses the configuration file securely.
 *
 * Implements robust file handling and input validation during parsing.
 * Uses std::unique_ptr internally to ensure memory is freed if parsing fails 
 * mid-way, preventing memory leaks.
 *
 * @param filename The path to the configuration file.
 * @return Configuration* or NULL.
 */
Configuration* ConfigLoader::loadConfig(const std::string& filename) {
    // Use unique_ptr for safe memory management internally until ownership is released.
    std::unique_ptr<Configuration> config_ptr = std::make_unique<Configuration>();
    
    std::ifstream file(filename);

    // SECURE PRACTICE: Check if the file opened successfully. 
    // This is the source of the potential NULL return.
    if (!file.is_open()) {
        std::cerr << "[ERROR] ConfigLoader: Failed to open configuration file: " << filename << std::endl;
        return nullptr; // Return NULL on critical failure
    }

    std::string line;
    int line_num = 0;

    while (std::getline(file, line)) {
        line_num++;
        
        // 1. Input Validation & Sanitization: Ignore comments and empty lines
        line = trim(line);
        if (line.empty() || line[0] == '#' || line[0] == ';') {
            continue;
        }

        // Find the delimiter '='
        size_t delimiter_pos = line.find('=');
        if (delimiter_pos == std::string::npos) {
            // SECURE PRACTICE: Log parsing errors and skip malformed lines.
            std::cerr << "[WARNING] ConfigLoader: Skipping malformed line " << line_num << " (no delimiter).\n";
            continue;
        }

        // Extract key and value
        std::string key = line.substr(0, delimiter_pos);
        std::string value = line.substr(delimiter_pos + 1);

        // 2. Input Validation: Trim whitespace from key and value
        key = trim(key);
        value = trim(value);

        // 3. Input Validation: Ensure key is not empty after trimming
        if (key.empty()) {
            std::cerr << "[WARNING] ConfigLoader: Empty key found on line " << line_num << ". Skipping.\n";
            continue;
        }

        // Add setting to the configuration object
        config_ptr->addSetting(key, value);
    }

    // Check if the file was read successfully (e.g., EOF reached without I/O errors)
    if (file.bad()) {
        std::cerr << "[ERROR] ConfigLoader: I/O error while reading file." << std::endl;
        return nullptr; // Return NULL on I/O error
    }

    // Transfer ownership of the Configuration object to the caller (raw pointer interface)
    return config_ptr.release();
}