#ifndef DATA_STORE_H
#define DATA_STORE_H

#include <stddef.h>

// Define the maximum number of configuration entries allowed.
#define MAX_CONFIG_ENTRIES 10
// Define the maximum size for the data payload (e.g., a short string or value).
#define MAX_DATA_SIZE 32

// Structure representing a single configuration entry
typedef struct {
    int id;
    char value[MAX_DATA_SIZE];
} ConfigEntry;

// Function prototypes for data management
/**
 * @brief Stores data at a specific index.
 * @param index The target index (must be < MAX_CONFIG_ENTRIES).
 * @param data The string data to store.
 * @return 0 on success, -1 on failure (out of bounds or invalid data).
 */
int store_data_at_index(size_t index, const char *data);

/**
 * @brief Initializes the data store.
 */
void initialize_data_store(void);

/**
 * @brief Displays the current contents of the data store.
 */
void display_data_store(void);

#endif // DATA_STORE_H