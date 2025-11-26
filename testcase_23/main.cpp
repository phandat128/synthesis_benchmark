#include <iostream>
#include <string>
#include <vector>
#include "ConfigParser.h"
#include "DataProcessor.h"

// Helper function to display usage instructions.
void print_usage(const std::string& program_name) {
    std::cout << "Usage: " << program_name << " --set-index <INDEX> --set-value <VALUE>" << std::endl;
    std::cout << "  <INDEX> must be an integer between 0 and " << BUFFER_SIZE - 1 << "." << std::endl;
    std::cout << "Example: " << program_name << " --set-index 5 --set-value 42" << std::endl;
}

int main(int argc, char* argv[]) {
    if (argc < 5) {
        print_usage(argv[0]);
        return 1;
    }

    std::string raw_index_input;
    std::string raw_value_input;
    bool index_found = false;
    bool value_found = false;

    // 1. Read raw command-line arguments (Source of Taint)
    for (int i = 1; i < argc; ++i) {
        std::string arg = argv[i];
        
        if (arg == "--set-index" && i + 1 < argc) {
            // Tainted input read
            raw_index_input = argv[i + 1];
            index_found = true;
            i++; // Skip the next argument (the value)
        } else if (arg == "--set-value" && i + 1 < argc) {
            raw_value_input = argv[i + 1];
            value_found = true;
            i++; // Skip the next argument (the value)
        }
    }

    if (!index_found || !value_found) {
        std::cerr << "Error: Both --set-index and --set-value must be provided." << std::endl;
        print_usage(argv[0]);
        return 1;
    }

    // Initialize components
    DataProcessor processor;
    processor.print_buffer();

    // 2. Parse and Validate Input (Propagation/Sanitization)
    // ConfigParser handles string conversion, exception handling, and crucial bounds checking.
    ConfigUpdate update = ConfigParser::parse_config(raw_index_input, raw_value_input);

    if (update.valid) {
        std::cout << "\n[MAIN] Configuration successfully parsed and validated." << std::endl;
        
        // 3. Apply Update (Sink is protected by DataProcessor's internal check)
        bool success = processor.update_setting(update.index, update.value);

        if (success) {
            processor.print_buffer();
            return 0;
        } else {
            // Defensive error handling for internal processor failure.
            std::cerr << "[MAIN ERROR] Failed to apply setting due to internal processor error." << std::endl;
            return 2;
        }

    } else {
        // Error handling: Input was invalid (non-numeric, out-of-bounds, etc.)
        std::cerr << "\n[MAIN ERROR] Configuration input failed validation. Aborting update." << std::endl;
        return 3;
    }
}