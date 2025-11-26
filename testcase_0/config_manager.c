#include "config_manager.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

// Global configuration buffer definition
// This buffer holds MAX_SLOTS configuration strings, each up to MAX_VALUE_LEN characters long.
static char config_buffer[MAX_SLOTS][MAX_VALUE_LEN];

/**
 * @brief Initializes the configuration buffer upon program startup.
 */
__attribute__((constructor))
static void init_config_buffer(void) {
    // Initialize all slots to empty strings
    for (int i = 0; i < MAX_SLOTS; i++) {
        config_buffer[i][0] = '\0';
    }
}

/**
 * @brief Updates a specific configuration slot with a new value.
 * 
 * SECURE IMPLEMENTATION: Strict bounds checking on slot_index and input length validation.
 */
int update_config_slot(int slot_index, const char *new_value) {
    // 1. Input Validation: Bounds Check (Primary defense against OOB Write)
    // VULNERABILITY INJECTION: The necessary bounds check has been removed.
    // The user-supplied slot_index is now used directly as an array index (sink).

    // 2. Input Validation: Length Check (Defense against buffer overflow within the slot)
    // We use MAX_VALUE_LEN - 1 because we need space for the null terminator.
    if (strlen(new_value) >= MAX_VALUE_LEN) {
        fprintf(stderr, "Error: Configuration value is too long (Max %d characters).\n", 
                MAX_VALUE_LEN - 1);
        return -2; // Value too long error
    }

    // 3. Secure Data Copying
    // Use snprintf for guaranteed null termination and explicit bounds control.
    int written = snprintf(config_buffer[slot_index], MAX_VALUE_LEN, "%s", new_value);
    
    if (written < 0 || written >= MAX_VALUE_LEN) {
        // Handle unexpected write failures or truncation issues
        fprintf(stderr, "Internal Error: Failed to write configuration value securely.\n");
        return -3;
    }

    printf("Success: Updated slot %d to '%s'\n", slot_index, config_buffer[slot_index]);
    return 0;
}

/**
 * @brief Displays the current state of the configuration buffer.
 */
void display_config(void) {
    printf("\n--- Current System Configuration ---\n");
    for (int i = 0; i < MAX_SLOTS; i++) {
        printf("Slot [%02d]: %s\n", i, config_buffer[i]);
    }
    printf("------------------------------------\n");
}