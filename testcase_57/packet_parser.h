#ifndef PACKET_PARSER_H
#define PACKET_PARSER_H

#include <stddef.h>

// Define a reasonable maximum size for security and resource management.
// This constant defines the actual allocated size of the packet buffer.
#define MAX_PACKET_SIZE 1024

// Define the maximum size we are willing to extract into a temporary buffer.
// This helps prevent potential buffer overflows in the destination buffer.
#define MAX_EXTRACT_SIZE 512

// Function prototypes

/**
 * @brief Securely extracts a segment of data from a source buffer.
 * 
 * This function implements critical bounds checking to prevent CWE-125 Out-of-bounds Read.
 * 
 * @param buffer The source buffer containing the packet data.
 * @param buffer_size The actual allocated size of the source buffer (e.g., MAX_PACKET_SIZE).
 * @param read_size The requested number of bytes to read (potentially tainted input).
 * @param destination The destination buffer where data will be copied.
 * @return The number of bytes successfully read, or a negative error code.
 */
int extract_data_segment(const char *buffer, size_t buffer_size, size_t read_size, char *destination);

/**
 * @brief Reads and validates an integer input from the user.
 * 
 * @param prompt The message displayed to the user.
 * @return The validated integer value, or -1 on error/invalid input.
 */
long read_integer_input(const char *prompt);

#endif // PACKET_PARSER_H