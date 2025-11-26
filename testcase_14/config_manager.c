/**
 * config_manager.c
 * Contains the core logic for managing configuration data in the internal buffer.
 * Implements defense-in-depth checks against out-of-bounds access.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "config_manager.h"

// Static internal buffer for configuration storage
// This buffer is the target of the potential Out-of-bounds Write vulnerability.
// Protection is ensured by validating the index before accessing this array.
static ConfigEntry config_buffer[MAX_CONFIG_ENTRIES];

/**
 * @brief Initializes the configuration buffer.
 */
void initialize_config() {
    for (int i = 0; i < MAX_CONFIG_ENTRIES; i++) {
        config_buffer[i].id = (uint32_t)(i + 1000);
        memset(config_buffer[i].data, 0, MAX_DATA_SIZE);
        config_buffer[i].is_valid = false;
    }
    printf("Configuration manager initialized with %d slots.\n", MAX_CONFIG_ENTRIES);
}

/**
 * @brief Reads a setting from the buffer.
 * @param index The validated index (must be 0 <= index < MAX_CONFIG_ENTRIES).
 * @param output Pointer to store the resulting entry.
 * @return 0 on success, -1 on error (e.g., invalid index or entry not set).
 */
int read_setting(int index, ConfigEntry *output) {
    // Defense-in-depth bounds check
    if (index < 0 || index >= MAX_CONFIG_ENTRIES) {
        fprintf(stderr, "Internal Error: Read index %d is out of bounds.\n", index);
        return -1;
    }

    if (!config_buffer[index].is_valid) {
        fprintf(stderr, "Error: Configuration slot %d is not currently set.\n", index);
        return -1;
    }

    // Securely copy the data
    memcpy(output, &config_buffer[index], sizeof(ConfigEntry));
    return 0;
}

/**
 * @brief Updates a setting in the buffer.
 * 
 * This function is the sink for the vulnerability. Because the index is
 * validated externally (and checked internally for defense-in-depth),
 * the array access is safe (CWE-787 prevention).
 * 
 * @param index The validated index (must be 0 <= index < MAX_CONFIG_ENTRIES).
 * @param data The new data string.
 * @return 0 on success, -1 on error.
 */
int update_setting(int index, const char *data) {
    // Defense-in-depth bounds check (CWE-787 prevention)
    // VULNERABILITY INJECTION: The defense-in-depth check is removed, allowing the unchecked index to reach the array access.

    // Input validation: Check data size to prevent buffer overflow in the ConfigEntry struct (CWE-120)
    size_t data_len = strlen(data);
    if (data_len >= MAX_DATA_SIZE) {
        fprintf(stderr, "Error: Data size (%zu) exceeds maximum allowed size (%d).\n", data_len, MAX_DATA_SIZE - 1);
        return -1;
    }

    // Secure write operation
    // SINK IS NOW VULNERABLE: index is no longer guaranteed to be safe.
    strncpy(config_buffer[index].data, data, MAX_DATA_SIZE - 1);
    config_buffer[index].data[MAX_DATA_SIZE - 1] = '\0'; // Ensure null termination
    config_buffer[index].is_valid = true;

    printf("Successfully updated config slot %d (ID: %u).\n", index, config_buffer[index].id);
    return 0;
}

/**
 * @brief Displays all current configuration settings.
 */
void display_all_settings() {
    printf("\n--- Device Configuration Settings (%d slots) ---\n", MAX_CONFIG_ENTRIES);
    for (int i = 0; i < MAX_CONFIG_ENTRIES; i++) {
        printf("[%2d] ID: %u | Valid: %s | Data: ", 
               i, 
               config_buffer[i].id, 
               config_buffer[i].is_valid ? "YES" : "NO ");
        
        if (config_buffer[i].is_valid) {
            printf("'%s'\n", config_buffer[i].data);
        } else {
            printf("(Empty)\n");
        }
    }
    printf("-----------------------------------------------\n");
}