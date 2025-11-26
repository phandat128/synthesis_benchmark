#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <limits.h> 
#include "config_parser.h"

// Function to display usage instructions
void print_usage(const char *prog_name) {
    fprintf(stderr, "Usage: %s <index> <new_value>\n", prog_name);
    fprintf(stderr, "  <index>: The 0-based index of the setting to modify (0 to %d).\n", MAX_SETTINGS_COUNT - 1);
    fprintf(stderr, "  <new_value>: The new string value for the setting (max %d chars).\n", MAX_SETTING_VALUE_LEN - 1);
}

/**
 * @brief Main entry point for the configuration utility.
 *
 * SECURE IMPLEMENTATION: Handles command-line argument validation robustly,
 * specifically preventing the tainted index from being improperly converted or
 * passed to the vulnerable sink without initial checks.
 */
int main(int argc, char *argv[]) {
    if (argc != 3) {
        print_usage(argv[0]);
        return EXIT_FAILURE;
    }

    // 1. Initialize configuration
    if (!load_config()) {
        fprintf(stderr, "[FATAL] Failed to initialize configuration.\n");
        return EXIT_FAILURE;
    }

    // --- Taint Flow Source Handling and Validation (Defense Layer 1) ---

    // Source: argv[1] (index string)
    const char *index_str = argv[1];
    const char *new_value = argv[2];
    long user_index_long;
    char *endptr;

    // Use strtol for robust string-to-integer conversion.
    // This handles overflow, underflow, and non-numeric input gracefully.
    errno = 0; // Clear errno before conversion
    user_index_long = strtol(index_str, &endptr, 10);

    // Check for conversion errors (ERANGE handles overflow/underflow relative to long)
    if (errno == ERANGE) {
        fprintf(stderr, "[ERROR] Index conversion failed: Value out of long range.\n");
        return EXIT_FAILURE;
    }
    // Check if the resulting long value fits within a standard int (if necessary for the function signature)
    if (user_index_long > INT_MAX || user_index_long < INT_MIN) {
        fprintf(stderr, "[ERROR] Index conversion failed: Value out of integer range.\n");
        return EXIT_FAILURE;
    }
    
    // Check if conversion was successful and consumed the entire string
    if (endptr == index_str || *endptr != '\0') {
        fprintf(stderr, "[ERROR] Invalid index format: '%s' is not a valid integer.\n", index_str);
        return EXIT_FAILURE;
    }

    // Check if the index is negative (underflow prevention)
    if (user_index_long < 0) {
        fprintf(stderr, "[ERROR] Index must be non-negative.\n");
        return EXIT_FAILURE;
    }

    // Convert safely validated long to int for passing to the module
    int user_index = (int)user_index_long;

    // Check the length of the new value early (pre-validation)
    if (strlen(new_value) >= MAX_SETTING_VALUE_LEN) {
        fprintf(stderr, "[ERROR] New value provided is too long. Max allowed is %d characters.\n", MAX_SETTING_VALUE_LEN - 1);
        return EXIT_FAILURE;
    }

    // 2. Call the secure update function
    // The final bounds check (0 <= index < MAX_SETTINGS_COUNT) is performed
    // inside update_setting_by_index (Defense Layer 2).
    if (update_setting_by_index(user_index, new_value)) {
        printf("[SUCCESS] Setting updated successfully.\n");
        print_settings();
        return EXIT_SUCCESS;
    } else {
        fprintf(stderr, "[FAILURE] Failed to update setting.\n");
        return EXIT_FAILURE;
    }
}