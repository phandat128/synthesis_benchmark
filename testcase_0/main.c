#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include "config_manager.h"

/**
 * @brief Main entry point for the configuration utility.
 * 
 * Usage: ./config_util <index> <value>
 * 
 * Implements robust argument parsing and validation before calling the core update logic.
 */
int main(int argc, char *argv[]) {
    long index_l;
    char *endptr;
    int slot_index;
    const char *new_value;

    if (argc != 3) {
        fprintf(stderr, "Usage: %s <slot_index> <new_value>\n", argv[0]);
        fprintf(stderr, "  <slot_index>: Integer between 0 and %d.\n", MAX_SLOTS - 1);
        fprintf(stderr, "  <new_value>: String value (Max length %d).\n", MAX_VALUE_LEN - 1);
        return EXIT_FAILURE;
    }

    // --- 1. Secure Parsing and Validation of Slot Index (argv[1]) ---

    // Reset errno before calling strtol to detect conversion errors reliably
    errno = 0; 
    
    // Use strtol for robust conversion, checking for overflow/underflow and invalid characters.
    index_l = strtol(argv[1], &endptr, 10);

    // Check for conversion errors (e.g., value too large/small for long)
    if (errno == ERANGE) {
        fprintf(stderr, "Error: Index value '%s' is out of range for a long integer.\n", argv[1]);
        return EXIT_FAILURE;
    }
    
    // Check if conversion stopped prematurely (i.e., non-numeric characters present)
    if (endptr == argv[1] || *endptr != '\0') {
        fprintf(stderr, "Error: Index value '%s' is not a valid integer.\n", argv[1]);
        return EXIT_FAILURE;
    }

    // Check if the parsed long fits safely into an int and is non-negative.
    // Note: The bounds check against MAX_SLOTS is performed inside update_config_slot.
    if (index_l < 0 || index_l > INT_MAX) {
        fprintf(stderr, "Error: Index value must be a non-negative integer and within system limits.\n");
        return EXIT_FAILURE;
    }
    
    // Cast safely to int after validation
    slot_index = (int)index_l;

    // --- 2. Retrieve Value (argv[2]) ---
    new_value = argv[2];

    // --- 3. Call Secure Update Function ---
    int result = update_config_slot(slot_index, new_value);

    if (result != 0) {
        // Error handling is robust and does not leak internal memory layout details.
        fprintf(stderr, "Configuration update failed (Error Code: %d). Check index and value length.\n", result);
        return EXIT_FAILURE;
    }

    // Display configuration for verification
    display_config();

    return EXIT_SUCCESS;
}