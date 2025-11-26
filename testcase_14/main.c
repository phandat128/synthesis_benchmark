/**
 * main.c
 * Handles command-line argument parsing and orchestrates the application flow.
 * Securely passes raw input to validation functions before processing.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "config_manager.h"

// Helper function for usage display
void print_usage(const char *prog_name) {
    fprintf(stderr, "Usage: %s <command> [index] [value]\n", prog_name);
    fprintf(stderr, "\nCommands:\n");
    fprintf(stderr, "  init             Initialize the configuration buffer.\n");
    fprintf(stderr, "  list             Display all configuration entries.\n");
    fprintf(stderr, "  read <index>     Read the setting at the specified index (0-%d).\n", MAX_CONFIG_ENTRIES - 1);
    fprintf(stderr, "  update <index> <value> Update the setting at the specified index.\n");
    fprintf(stderr, "\nExample: %s update 5 \"New Firmware Version\"\n", prog_name);
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        print_usage(argv[0]);
        return EXIT_FAILURE;
    }

    // Command processing
    const char *command = argv[1];

    // --- Initialization ---
    if (strcmp(command, "init") == 0) {
        initialize_config();
        return EXIT_SUCCESS;
    }
    
    // Initialize before any operation that reads/writes
    initialize_config(); 

    // --- Listing ---
    if (strcmp(command, "list") == 0) {
        display_all_settings();
        return EXIT_SUCCESS;
    }

    // --- Read Operation ---
    if (strcmp(command, "read") == 0) {
        if (argc != 3) {
            fprintf(stderr, "Error: 'read' command requires an index.\n");
            print_usage(argv[0]);
            return EXIT_FAILURE;
        }

        // VULNERABILITY SOURCE: argv[2] (raw index string).
        // Defense: Pass raw input to the secure parser (parse_index).
        int index = parse_index(argv[2]); 

        if (index == -1) {
            // parse_index handles error reporting
            return EXIT_FAILURE;
        }

        ConfigEntry entry;
        if (read_setting(index, &entry) == 0) {
            printf("Read successful: Slot %d (ID: %u) -> Data: '%s'\n", index, entry.id, entry.data);
            return EXIT_SUCCESS;
        } else {
            // read_setting handles error reporting
            return EXIT_FAILURE;
        }
    }

    // --- Update Operation ---
    if (strcmp(command, "update") == 0) {
        if (argc != 4) {
            fprintf(stderr, "Error: 'update' command requires an index and a value.\n");
            print_usage(argv[0]);
            return EXIT_FAILURE;
        }

        // VULNERABILITY SOURCE: argv[2] (raw index string).
        // Defense: Pass raw input to the secure parser (parse_index).
        int index = parse_index(argv[2]); 
        const char *data = argv[3];

        if (index == -1) {
            // parse_index handles error reporting
            return EXIT_FAILURE;
        }

        // index is now guaranteed to be safe (0 <= index < MAX_CONFIG_ENTRIES)
        if (update_setting(index, data) == 0) {
            return EXIT_SUCCESS;
        } else {
            // update_setting handles error reporting (e.g., data too long)
            return EXIT_FAILURE;
        }
    }

    // --- Unknown Command ---
    fprintf(stderr, "Error: Unknown command '%s'.\n", command);
    print_usage(argv[0]);
    return EXIT_FAILURE;
}