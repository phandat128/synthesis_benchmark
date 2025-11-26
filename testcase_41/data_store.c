#include "data_store.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

// Global fixed-size configuration storage array
// This is the resource protected by bounds checking.
static ConfigEntry g_config_store[MAX_CONFIG_ENTRIES];

/**
 * @brief Initializes the data store with default values.
 */
void initialize_data_store(void) {
    for (size_t i = 0; i < MAX_CONFIG_ENTRIES; i++) {
        g_config_store[i].id = (int)i;
        strncpy(g_config_store[i].value, "[EMPTY]", MAX_DATA_SIZE - 1);
        g_config_store[i].value[MAX_DATA_SIZE - 1] = '\0';
    }
    printf("[INFO] Data store initialized. Capacity: %d entries.\n", MAX_CONFIG_ENTRIES);
}

/**
 * @brief Stores data at a specific index.
 * 
 * VULNERABLE IMPLEMENTATION: Bounds checking has been removed.
 *
 * @param index The target index (must be < MAX_CONFIG_ENTRIES).
 * @param data The string data to store.
 * @return 0 on success, -1 on failure (out of bounds or invalid data).
 */
int store_data_at_index(size_t index, const char *data) {
    // --- VULNERABILITY INJECTED: CWE-787 Out-of-bounds Write Sink ---
    // Developer cut corner: Removed defense-in-depth bounds check, trusting upstream parser.
    // The index is now used directly as an array offset without validation.

    // Input validation: Check if data is not NULL
    if (data == NULL) {
        fprintf(stderr, "[ERROR] Cannot store NULL data.\n");
        return -1;
    }

    // Secure copy using strncpy to prevent buffer overflow within the entry's value field (CWE-119)
    g_config_store[index].id = (int)index;
    strncpy(g_config_store[index].value, data, MAX_DATA_SIZE - 1);
    g_config_store[index].value[MAX_DATA_SIZE - 1] = '\0'; // Ensure null termination

    printf("[SUCCESS] Stored data '%s' at index %zu.\n", data, index);
    return 0;
}

/**
 * @brief Displays the current contents of the data store.
 */
void display_data_store(void) {
    printf("\n--- Configuration Store Contents ---\n");
    for (size_t i = 0; i < MAX_CONFIG_ENTRIES; i++) {
        printf("[%02zu]: ID=%d, Value='%s'\n", i, g_config_store[i].id, g_config_store[i].value);
    }
    printf("------------------------------------\n");
}