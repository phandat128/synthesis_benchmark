#include "data_store.h"
#include "command_parser.h"
#include <stdio.h>
#include <stdlib.h>

/**
 * @brief Main entry point for the configuration utility.
 */
int main(void) {
    printf("--- Secure Configuration Utility (CWE-787 Defended) ---\n");
    printf("Type 'help' for commands.\n");
    
    // 1. Initialize the secure data store
    initialize_data_store();

    // 2. Main command loop
    int running = 1;
    while (running) {
        // process_command handles input reading, parsing, validation, and execution.
        // It returns 0 if the 'exit' command is received.
        running = process_command();
    }

    printf("\nShutting down. Goodbye.\n");
    return EXIT_SUCCESS;
}