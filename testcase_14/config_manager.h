/**
 * config_manager.h
 * Defines the structure for configuration entries and function prototypes.
 */

#ifndef CONFIG_MANAGER_H
#define CONFIG_MANAGER_H

#include <stdint.h>
#include <stdbool.h>
#include <limits.h>

// Security Note: Define a fixed, reasonable size for the configuration buffer.
#define MAX_CONFIG_ENTRIES 16
#define MAX_DATA_SIZE 64

// Define the structure for a single configuration entry
typedef struct {
    uint32_t id;
    char data[MAX_DATA_SIZE];
    bool is_valid;
} ConfigEntry;

// Function Prototypes for Configuration Management
void initialize_config();
int read_setting(int index, ConfigEntry *output);
int update_setting(int index, const char *data);
void display_all_settings();

// Function Prototypes for Input Parsing (Securely handles index validation)
// Returns the validated index (>= 0) or -1 on failure (out of bounds or invalid format).
int parse_index(const char *input);

#endif // CONFIG_MANAGER_H