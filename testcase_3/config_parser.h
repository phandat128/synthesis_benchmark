#ifndef CONFIG_PARSER_H
#define CONFIG_PARSER_H

#include <stdbool.h>

// --- Security Constants ---
// Define the maximum number of settings allowed. This is the critical boundary
// that must be checked against user input indices (CWE-787 Defense).
#define MAX_SETTINGS_COUNT 32
// Define the maximum length for a setting value (CWE-120 Defense).
#define MAX_SETTING_VALUE_LEN 64

// Structure to hold a single configuration setting
typedef struct {
    char key[MAX_SETTING_VALUE_LEN];
    char value[MAX_SETTING_VALUE_LEN];
    bool is_valid;
} config_setting_t;

// Function Prototypes
/**
 * @brief Loads configuration settings into the internal buffer.
 * @return true on success, false otherwise.
 */
bool load_config(void);

/**
 * @brief Updates a configuration setting based on its index.
 *
 * This function implements the crucial bounds checking to prevent CWE-787.
 *
 * @param index The 0-based index of the setting to update.
 * @param new_value The new string value for the setting.
 * @return true on successful update, false if index is out of bounds or value is too long.
 */
bool update_setting_by_index(int index, const char *new_value);

/**
 * @brief Prints all currently loaded settings.
 */
void print_settings(void);

#endif // CONFIG_PARSER_H