#include <iostream>
#include <string>
#include <limits>
#include "CommandProcessor.h"
#include "DataBuffer.h"

/**
 * @brief Displays usage instructions to the user.
 */
void display_help() {
    std::cout << "\n--- Secure Configuration Utility ---\n";
    std::cout << "Buffer Size: " << BUFFER_SIZE << " elements (Indices 0 to " << BUFFER_SIZE - 1 << ")\n";
    std::cout << "Commands:\n";
    std::cout << "  SET <value> <index> : Write <value> to <index>. (e.g., SET 100 5)\n";
    std::cout << "  GET <index>         : Read value from <index>. (e.g., GET 5)\n";
    std::cout << "  EXIT / QUIT         : Terminate the program.\n";
    std::cout << "------------------------------------\n";
}

int main() {
    // 1. Initialize secure components
    DataBuffer data_store;
    CommandProcessor processor(data_store);

    display_help();

    std::string input_line;
    bool running = true;

    // 2. Main command loop
    while (running) {
        std::cout << "\n> ";
        
        // Secure Input Handling: Read entire line from standard input.
        if (!std::getline(std::cin, input_line)) {
            // Handle EOF or read error gracefully
            std::cerr << "\n[INFO] End of input detected." << std::endl;
            running = false;
            break;
        }

        // 3. Process the command. The processor handles all validation and security checks.
        running = processor.process_command(input_line);
        
        // Ensure error messages are flushed immediately
        std::cerr.flush();
    }

    std::cout << "Application terminated securely." << std::endl;
    return 0;
}