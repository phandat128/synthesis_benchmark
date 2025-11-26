#include "packet_parser.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

// Define a safe buffer size for reading input lines
#define INPUT_LINE_MAX 128

/**
 * @brief Reads and validates an integer input from the user.
 * 
 * Implements robust input handling using fgets and strtol to prevent buffer overflows
 * and correctly handle non-numeric input, ensuring the output is a positive long.
 * 
 * @param prompt The message displayed to the user.
 * @return The validated integer value, or -1 on error/invalid input.
 */
long read_integer_input(const char *prompt) {
    char line[INPUT_LINE_MAX];
    char *endptr;
    long val;

    printf("%s", prompt);

    // Use fgets for safe input reading (prevents buffer overflow on input line)
    if (fgets(line, sizeof(line), stdin) == NULL) {
        fprintf(stderr, "Error reading input.\n");
        return -1;
    }

    // Check if the input line was too long (buffer overflow attempt)
    if (line[strlen(line) - 1] != '\n' && !feof(stdin)) {
        fprintf(stderr, "Input too long or invalid format. Please try again.\n");
        // Clear the rest of the input buffer
        int c;
        while ((c = getchar()) != '\n' && c != EOF);
        return -1;
    }

    // Convert string to long integer
    errno = 0; // Reset errno before strtol
    val = strtol(line, &endptr, 10);

    // Check for conversion errors (overflow/underflow)
    if (errno == ERANGE) {
        fprintf(stderr, "Input value out of range.\n");
        return -1;
    }

    // Check for non-numeric input or empty input
    if (endptr == line || (*endptr != '\0' && *endptr != '\n')) {
        fprintf(stderr, "Invalid input: Please enter a numeric value.\n");
        return -1;
    }

    // Security check: Ensure the length is positive
    if (val <= 0) {
        fprintf(stderr, "Input must be a positive integer.\n")
        return -1;
    }

    return val;
}