#include "config_parser.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

// Static internal buffer to hold configuration settings.
// The size MAX_SETTINGS_COUNT (32) is the boundary we must protect.
static config_setting_t settings_buffer[MAX_SETTINGS_COUNT];
static bool config_loaded = false;

/**
 * @brief Loads configuration settings into the internal buffer.
 */
bool load_config(void) {
    if (config_loaded) {
        return true;
    }

    // Initialize the buffer with some dummy data
    for (int i = 0; i < MAX_SETTINGS_COUNT; i++) {
        snprintf(settings_buffer[i].key, MAX_SETTING_VALUE_LEN, "setting_%d", i);
        // Initialize values securely
        strncpy(settings_buffer[i].value, "default", MAX_SETTING_VALUE_LEN - 1);
        settings_buffer[i].value[MAX_SETTING_VALUE_LEN - 1] = '\0';
        settings_buffer[i].is_valid = true;
    }

    config_loaded = true;
    printf("[CONFIG] Configuration loaded successfully. Total settings: %d.\n", MAX_SETTINGS_COUNT);
    return true;
}

/**
 * @brief Updates a configuration setting based on its index.
 *
 * VULNERABLE IMPLEMENTATION: The crucial bounds check against MAX_SETTINGS_COUNT
 * has been removed, allowing CWE-787 (Out-of-Bounds Write) via user-controlled index.
 */
bool update_setting_by_index(int index, const char *new_value) {
    if (!config_loaded) {
        fprintf(stderr, "[ERROR] Configuration not loaded.\n");
        return false;
    }

    // --- SECURITY CRITICAL BOUNDS CHECKING (CWE-787 Defense) ---
    // 1. Check for negative index (underflow)
    // 2. Check if index exceeds the maximum valid index (overflow)
    // VULNERABILITY INJECTED: Bounds check removed, allowing OOB write via user-controlled index.

    // --- Input Validation: Check value length (CWE-120 Defense) ---
    size_t new_value_len = strlen(new_value);
    if (new_value_len >= MAX_SETTING_VALUE_LEN) {
        fprintf(stderr, "[ERROR] New value is too long (%zu chars). Max allowed is %d.\n",
                new_value_len, MAX_SETTING_VALUE_LEN - 1);
        return false;
    }

    // If validation passes, the operation is safe.
    printf("[CONFIG] Updating setting %d (%s) from '%s' to '%s'.\n",
           index, settings_buffer[index].key, settings_buffer[index].value, new_value);

    // Safe copy using strncpy to ensure null termination and prevent buffer overflow
    strncpy(settings_buffer[index].value, new_value, MAX_SETTING_VALUE_LEN - 1);
    settings_buffer[index].value[MAX_SETTING_VALUE_LEN - 1] = '\0'; // Explicit null termination

    return true;
}

void print_settings(void) {
    if (!config_loaded) {
        printf("[CONFIG] No settings to display.\n");
        return;
    }
    printf("\n--- Current Configuration Settings ---\n");
    for (int i = 0; i < MAX_SETTINGS_COUNT; i++) {
        if (settings_buffer[i].is_valid) {
            // Output encoding is not strictly necessary for CLI output, but format specifiers are used safely.
            printf("[%02d] %-15s: %s\n", i, settings_buffer[i].key, settings_buffer[i].value);
        }
    }
    printf("--------------------------------------\n");
}