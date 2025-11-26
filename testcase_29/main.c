#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <limits.h>
#include "data_processor.h"
#include "config.h"

/**
 * @brief Safely converts a string to a long integer, validating bounds and format.
 * @param str The string to convert.
 * @param val Pointer to store the result.
 * @return 0 on success, -1 on failure.
 */
static int safe_strtol(const char *str, long *val) {
    char *endptr;
    long temp_val;

    // Check for null or empty string
    if (str == NULL || *str == '\0') {
        return -1;
    }

    errno = 0; // Clear errno before call
    temp_val = strtol(str, &endptr, 10);

    // Check for conversion errors (ERANGE, invalid characters)
    if ((errno == ERANGE && (temp_val == LONG_MAX || temp_val == LONG_MIN)) || (errno != 0 && temp_val == 0)) {
        fprintf(stderr, "[ERROR] Conversion error or value out of long range.\n");
        return -1;
    }

    // Check if the entire string was consumed (no trailing non-numeric characters)
    if (endptr == str || *endptr != '\0') {
        fprintf(stderr, "[ERROR] Invalid numeric format: '%s'.\n", str);
        return -1;
    }

    *val = temp_val;
    return 0;
}

/**
 * @brief Main entry point for the configuration utility.
 * Usage: ./config_util <index> <length>
 */
int main(int argc, char *argv[]) {
    if (argc != 3) {
        fprintf(stderr, "Usage: %s <segment_index> <segment_length>\n", argv[0]);
        fprintf(stderr, "  <segment_index>: Starting byte index (0 to %d)\n", MAX_VALID_INDEX);
        fprintf(stderr, "  <segment_length>: Number of bytes to read (1 to %d)\n", MAX_SEGMENT_SIZE);
        return EXIT_FAILURE;
    }

    long requested_index_l, requested_length_l;
    int requested_index, requested_length;

    // 1. Input Validation: Index (argv[1])
    if (safe_strtol(argv[1], &requested_index_l) != 0) {
        fprintf(stderr, "[FATAL] Invalid index argument provided.\n");
        return EXIT_FAILURE;
    }

    // 2. Input Validation: Length (argv[2] - The vulnerability source)
    if (safe_strtol(argv[2], &requested_length_l) != 0) {
        fprintf(stderr, "[FATAL] Invalid length argument provided.\n");
        return EXIT_FAILURE;
    }

    // 3. Range and Type Casting Check
    // Ensure values fit safely into standard 'int' type used by data_processor functions
    if (requested_index_l < 0 || requested_index_l > INT_MAX ||
        requested_length_l < 0 || requested_length_l > INT_MAX) {
        fprintf(stderr, "[FATAL] Index or length value is negative or too large for internal processing (int overflow).\n");
        return EXIT_FAILURE;
    }

    requested_index = (int)requested_index_l;
    requested_length = (int)requested_length_l;

    // 4. Application-Specific Boundary Checks (Initial sanity check)
    // Note: data_processor.c performs the final, definitive bounds check against MAX_BUFFER_SIZE.
    if (requested_index > MAX_VALID_INDEX) {
        fprintf(stderr, "[FATAL] Index %d is outside the valid range [0, %d].\n", requested_index, MAX_VALID_INDEX);
        return EXIT_FAILURE;
    }

    if (requested_length <= 0 || requested_length > MAX_SEGMENT_SIZE) {
        fprintf(stderr, "[FATAL] Length %d is outside the valid range [1, %d].\n", requested_length, MAX_SEGMENT_SIZE);
        return EXIT_FAILURE;
    }

    // Initialize the data buffer
    initialize_buffer();

    // 5. Process the request
    SegmentResult result = extract_data_segment(requested_index, requested_length);

    if (result.success) {
        printf("\n--- Extracted Data Segment ---\n");
        printf("Start Index: %d\n", requested_index);
        printf("Actual Length: %zu bytes\n", result.length);
        printf("Data (Hex): ");
        for (size_t i = 0; i < result.length; i++) {
            printf("%02X ", (unsigned char)result.data[i]);
        }
        printf("\nData (ASCII): ");
        for (size_t i = 0; i < result.length; i++) {
            // Output Encoding: Only print printable characters
            if (result.data[i] >= 32 && result.data[i] <= 126) {
                printf("%c", result.data[i]);
            } else {
                printf(".");
            }
        }
        printf("\n------------------------------\n");
        return EXIT_SUCCESS;
    } else {
        // Proper Error Handling: Do not leak internal state, rely on specific error logs from data_processor.c
        fprintf(stderr, "[FATAL] Data extraction failed due to invalid parameters or internal error.\n");
        return EXIT_FAILURE;
    }
}