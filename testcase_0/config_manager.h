#ifndef CONFIG_MANAGER_H
#define CONFIG_MANAGER_H

#include <limits.h>

// Define constants for the configuration system
#define MAX_SLOTS 10        // Maximum number of configuration parameters
#define MAX_VALUE_LEN 64    // Maximum length for a configuration value string (including null terminator)

// Function prototypes
/**
 * @brief Updates a specific configuration slot with a new value.
 * 
 * Implements strict bounds checking to prevent Out-of-Bounds Write vulnerabilities.
 * 
 * @param slot_index The index of the configuration slot to update (0 to MAX_SLOTS - 1).
 * @param new_value The new string value to store.
 * @return 0 on success, -1 if index is out of bounds, -2 if value is too long, -3 on internal write error.
 */
int update_config_slot(int slot_index, const char *new_value);

/**
 * @brief Displays the current state of the configuration buffer.
 */
void display_config(void);

#endif // CONFIG_MANAGER_H