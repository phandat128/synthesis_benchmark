#ifndef CONFIGPARSER_H
#define CONFIGPARSER_H

#include <string>
#include <optional>
#include "DataProcessor.h" // Need BUFFER_SIZE constant

// Structure to hold parsed and validated configuration data.
struct ConfigUpdate {
    int index = -1;
    int value = 0;
    bool valid = false;
};

class ConfigParser {
public:
    // Parses and validates the index and value strings.
    // Returns a ConfigUpdate struct indicating success and containing validated data.
    static ConfigUpdate parse_config(const std::string& raw_index, const std::string& raw_value);

private:
    // Helper function for secure string-to-integer conversion and bounds checking.
    static std::optional<int> parse_and_validate_index(const std::string& input);
    static std::optional<int> parse_value(const std::string& input);
};

#endif // CONFIGPARSER_H