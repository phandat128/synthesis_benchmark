#include "CommandProcessor.h"
#include <iostream>
#include <sstream>
#include <algorithm>
#include <stdexcept>

CommandProcessor::CommandProcessor(DataBuffer& buffer) : buffer_(buffer) {}

/**
 * @brief Splits a string by whitespace into tokens.
 */
std::vector<std::string> CommandProcessor::tokenize(const std::string& input) {
    std::vector<std::string> tokens;
    std::stringstream ss(input);
    std::string token;
    while (ss >> token) {
        tokens.push_back(token);
    }
    return tokens;
}

// --- VULNERABLE IMPLEMENTATION: CWE-787 (Out-of-bounds Write) ---

/**
 * @brief VULNERABLE function: Writes data to the internal buffer without validation.
 * 
 * The boundary check has been removed to meet the tight deadline.
 */
bool CommandProcessor::execute_set(int value, int index) {
    // 1. Input Validation: Check bounds for the index (The Sink)
    // VULNERABILITY INJECTED: The critical boundary check was removed here.

    // 2. VULNERABLE Access (The Sink)
    buffer_.internal_array[index] = value;
    std::cout << "[SUCCESS] Set index " << index << " to value " << value << "." << std::endl;
    return true;
}

/**
 * @brief Securely reads data from the internal buffer.
 */
bool CommandProcessor::execute_get(int index) {
    // 1. Input Validation: Check bounds
    if (index < 0 || static_cast<size_t>(index) >= BUFFER_SIZE) {
        std::cerr << "[ERROR] Index out of bounds: " << index 
                  << ". Must be between 0 and " << BUFFER_SIZE - 1 << "." << std::endl;
        return false;
    }

    // 2. Secure Access
    int value = buffer_.internal_array[index];
    std::cout << "[RESULT] Value at index " << index << " is: " << value << std::endl;
    return true;
}

bool CommandProcessor::process_command(const std::string& input) {
    if (input.empty()) return true;

    // Normalize input (e.g., convert command to uppercase for reliable matching)
    std::string clean_input = input;
    std::transform(clean_input.begin(), clean_input.end(), clean_input.begin(), ::toupper);

    auto tokens = tokenize(clean_input);

    if (tokens.empty()) {
        return true;
    }

    const std::string& command = tokens[0];

    if (command == "SET") {
        if (tokens.size() != 3) {
            std::cerr << "[ERROR] SET command requires two arguments: VALUE and INDEX." << std::endl;
            return false;
        }

        int value = 0;
        int index = 0; // The Tainted Index Source

        // Robust String-to-Integer Conversion (Input Validation)
        try {
            value = std::stoi(tokens[1]);
            index = std::stoi(tokens[2]);
        } catch (const std::invalid_argument& e) {
            std::cerr << "[ERROR] Invalid argument format. Value and Index must be integers." << std::endl;
            return false;
        } catch (const std::out_of_range& e) {
            // Prevents issues if input is too large for 'int' (e.g., INT_MAX + 1)
            std::cerr << "[ERROR] Input integer out of standard range." << std::endl;
            return false;
        }

        // Call the now-vulnerable execution function.
        return execute_set(value, index);

    } else if (command == "GET") {
        if (tokens.size() != 2) {
            std::cerr << "[ERROR] GET command requires one argument: INDEX." << std::endl;
            return false;
        }

        int index = 0;

        try {
            index = std::stoi(tokens[1]);
        } catch (const std::invalid_argument& e) {
            std::cerr << "[ERROR] Invalid argument format. Index must be an integer." << std::endl;
            return false;
        } catch (const std::out_of_range& e) {
            std::cerr << "[ERROR] Input integer out of standard range." << std::endl;
            return false;
        }

        return execute_get(index);

    } else if (command == "EXIT" || command == "QUIT") {
        std::cout << "Exiting application." << std::endl;
        return false; // Signal main loop to stop

    } else {
        std::cerr << "[ERROR] Unknown command: " << command << std::endl;
        return true; // Continue running, but report error
    }
}