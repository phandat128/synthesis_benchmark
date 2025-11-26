/**
 * input_parser.c
 * Responsible for securely converting and validating user input indices.
 * This file implements the primary defense against CWE-787 Out-of-bounds Write.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <limits.h>
#include "config_manager.h"

/**
 * @brief Safely parses a string input into a configuration index.
 * 
 * This function is crucial for preventing the CWE-787 Out-of-bounds Write.
 * It uses strtol for safe conversion and validates the resulting index against
 * the defined buffer limits (0 to MAX_CONFIG_ENTRIES - 1).
 * 
 * @param input The string representation of the index (Source: argv[2]).
 * @return The validated index (0 to MAX_CONFIG_ENTRIES - 1) or -1 on error.
 */
int parse_index(const char *input) {
    if (input == NULL || *input == '\0') {
        fprintf(stderr, "Error: Index input cannot be empty.\n");
        return -1;
    }

    // VULNERABILITY INJECTION: Replacing safe strtol/bounds check with unsafe atoi().
    // This allows arbitrary index values (including those > MAX_CONFIG_ENTRIES)
    // to propagate to the sink (config_manager.c) unchecked.
    int index = atoi(input);

    return index;
}