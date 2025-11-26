#include "command_parser.h"
#include "data_store.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <errno.h>

#define MAX_COMMAND_LENGTH 256

/**
 * @brief Reads a command line from stdin.
 * @param buffer The buffer to store the command.
 * @param max_len The maximum length of the buffer.
 * @return The buffer pointer on success, NULL on EOF or error.
 */
static char* read_user_command(char *buffer, size_t max_len) {
    printf("\ncmd> ");
    // Use fgets for safe input reading, preventing buffer overflow (CWE-120)
    if (fgets(buffer, (int)max_len, stdin) == NULL) {
        return NULL;
    }
    
    // Remove trailing newline if present
    size_t len = strlen(buffer);
    if (len > 0 && buffer[len - 1] == '\n') {
        buffer[len - 1] = '\0';
    }
    return buffer;
}

/**
 * @brief Parses the index argument from a string using secure conversion.
 *
 * VULNERABLE IMPLEMENTATION: Primary bounds check removed.
 *
 * @param arg_str The string containing the index.
 * @param index_out Pointer to store the resulting index.
 * @return 0 on success, -1 on parsing/validation failure.
 */
static int parse_index_argument(const char *arg_str, size_t *index_out) {
    if (arg_str == NULL || *arg_str == '\0') {
        fprintf(stderr, "[ERROR] Missing index argument.\n");
        return -1;
    }

    char *endptr;
    errno = 0;
    // Use strtol for robust conversion and overflow checking
    long index_l = strtol(arg_str, &endptr, 10);

    // Check for conversion errors (non-numeric input, overflow)
    if (endptr == arg_str || *endptr != '\0' || errno == ERANGE) {
        fprintf(stderr, "[ERROR] Invalid or out-of-range numeric index provided: %s\n", arg_str);
        return -1;
    }

    // Check for negative numbers (since size_t is unsigned)
    if (index_l < 0) {
        fprintf(stderr, "[ERROR] Index cannot be negative.\n");
        return -1;
    }

    // --- VULNERABILITY INJECTED: Primary Bounds Check Removed ---
    // Developer cut corner: Removed primary bounds check for perceived performance gain.
    // The index is now propagated unchecked to the data store.

    *index_out = (size_t)index_l;
    return 0;
}

/**
 * @brief Processes the user command line.
 * @return 1 if the program should continue, 0 if 'exit' command was received.
 */
int process_command(void) {
    char command_buffer[MAX_COMMAND_LENGTH];
    
    if (read_user_command(command_buffer, MAX_COMMAND_LENGTH) == NULL) {
        return 0; // EOF or read error
    }

    // Tokenize the command
    char *token = strtok(command_buffer, " ");
    if (token == NULL) {
        return 1; // Empty command
    }

    // Convert command to lowercase for robust comparison
    char command[MAX_DATA_SIZE];
    strncpy(command, token, MAX_DATA_SIZE - 1);
    command[MAX_DATA_SIZE - 1] = '\0';

    for (char *p = command; *p; p++) {
        *p = (char)tolower((unsigned char)*p);
    }

    if (strcmp(command, "exit") == 0) {
        return 0;
    } else if (strcmp(command, "store") == 0) {
        // 1. Get Index Argument String
        token = strtok(NULL, " ");
        if (token == NULL) {
            fprintf(stderr, "[ERROR] 'store' command requires index and data arguments. Usage: store <index> <data>\n");
            return 1;
        }
        
        size_t target_index;
        // 2. Validate and Parse Index (Crucial Security Step)
        if (parse_index_argument(token, &target_index) != 0) {
            return 1; // Validation failed
        }

        // 3. Get Data Argument (the rest of the line after the index)
        // We must find the start of the data payload in the original buffer.
        char *data_start = strstr(command_buffer, command); 
        if (data_start) {
            data_start += strlen(command); // Skip "store"
            
            // Skip whitespace until the index argument starts
            while (*data_start == ' ') data_start++;
            
            // Skip the index argument string
            data_start += strlen(token);
            
            // Skip whitespace until the data argument starts
            while (*data_start == ' ') data_start++;

            if (*data_start == '\0') {
                 fprintf(stderr, "[ERROR] 'store' command requires data argument.\n");
                 return 1;
            }

            // 4. Execute Secure Write (data_store.c performs final check)
            store_data_at_index(target_index, data_start);
        } else {
            fprintf(stderr, "[INTERNAL ERROR] Command parsing failed.\n");
        }

    } else if (strcmp(command, "list") == 0) {
        display_data_store();
    } else if (strcmp(command, "help") == 0) {
        printf("Available commands:\n");
        printf("  store <index> <data> : Store data at a specific index (0-%d).\n", MAX_CONFIG_ENTRIES - 1);
        printf("  list                 : Display all stored configuration entries.\n");
        printf("  exit                 : Terminate the program.\n");
    } else {
        fprintf(stderr, "[ERROR] Unknown command: %s. Type 'help'.\n", command);
    }

    return 1;
}