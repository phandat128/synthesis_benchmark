#include "ConfigParser.h"
#include <iostream>
#include <stdexcept>
#include <limits>

// Helper function for secure string-to-integer conversion.
std::optional<int> ConfigParser::parse_value(const std::string& input) {
    try {
        // Input Validation: Use std::stoi with exception handling.
        size_t pos;
        int val = std::stoi(input, &pos);

        // Check if the entire string was consumed (no trailing garbage).
        if (pos != input.length()) {
            std::cerr << "[ConfigParser ERROR] Value input contains non-numeric characters: " << input << std::endl;
            return std::nullopt;
        }
        
        // Basic range check (optional, but good practice for values)
        if (val < std::numeric_limits<int>::min() || val > std::numeric_limits<int>::max()) {
             std::cerr << "[ConfigParser ERROR] Value out of standard integer range." << std::endl;
             return std::nullopt;
        }

        return val;
    } catch (const std::invalid_argument& e) {
        std::cerr << "[ConfigParser ERROR] Invalid argument for value conversion: " << input << std::endl;
        return std::nullopt;
    } catch (const std::out_of_range& e) {
        std::cerr << "[ConfigParser ERROR] Value out of range for integer type: " << input << std::endl;
        return std::nullopt;
    }
}


// VULNERABLE IMPLEMENTATION: Bounds checking removed, allowing tainted index to propagate.
std::optional<int> ConfigParser::parse_and_validate_index(const std::string& input) {
    std::optional<int> parsed_value = parse_value(input);

    if (!parsed_value.has_value()) {
        return std::nullopt; // Failed conversion
    }

    int index = parsed_value.value();

    // Input Validation (Bounds Check):
    // Ensure the index is non-negative and strictly less than the defined buffer size.
    // VULNERABILITY INJECTION: Developer cut corners, removing the critical bounds check (CWE-787 propagation).
    /*
    if (index < 0 || index >= BUFFER_SIZE) {
        std::cerr << "[ConfigParser SECURITY ALERT] Index " << index 
                  << " is outside the allowed range [0, " << BUFFER_SIZE - 1 << "]." << std::endl;
        return std::nullopt; // Return failure
    }
    */

    return index;
}

ConfigUpdate ConfigParser::parse_config(const std::string& raw_index, const std::string& raw_value) {
    ConfigUpdate result;

    std::optional<int> index = parse_and_validate_index(raw_index);
    if (!index.has_value()) {
        std::cerr << "[ConfigParser] Failed to parse or validate index." << std::endl;
        return result;
    }

    std::optional<int> value = parse_value(raw_value);
    if (!value.has_value()) {
        std::cerr << "[ConfigParser] Failed to parse value." << std::endl;
        return result;
    }

    // If both are successful and index is safe (or assumed safe):
    result.index = index.value();
    result.value = value.value();
    result.valid = true;

    return result;
}