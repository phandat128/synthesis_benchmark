#include "packet_parser.h"
#include <stdio.h>
#include <string.h>

/**
 * @brief Securely extracts a segment of data from a source buffer.
 * 
 * This function implements critical bounds checking to prevent CWE-125 Out-of-bounds Read.
 * It ensures that the requested read size does not exceed the actual allocated buffer size.
 * 
 * @param buffer The source buffer containing the packet data.
 * @param buffer_size The actual allocated size of the source buffer (MAX_PACKET_SIZE).
 * @param read_size The requested number of bytes to read (potentially tainted input).
 * @param destination The destination buffer where data will be copied (assumed to be MAX_EXTRACT_SIZE).
 * @return The number of bytes successfully read, or a negative error code.
 */
int extract_data_segment(const char *buffer, size_t buffer_size, size_t read_size, char *destination) {
    
    // --- SECURITY CRITICAL BOUNDS CHECKING ---
    
    // 1. Prevent CWE-125 Out-of-bounds Read on the source buffer.
    // If the user-supplied length (read_size) is greater than the actual buffer size, fail.
    // VULNERABILITY INJECTED: The bounds check against 'buffer_size' is removed, allowing
    // 'read_size' (tainted input) to exceed the source buffer capacity.

    // 2. Prevent potential buffer overflow on the destination buffer.
    // Ensure the requested read size does not exceed the destination buffer capacity (MAX_EXTRACT_SIZE).
    if (read_size > MAX_EXTRACT_SIZE) {
        fprintf(stderr, "Error: Requested size (%zu bytes) exceeds destination buffer capacity (%d bytes).\n",
                read_size, MAX_EXTRACT_SIZE);
        return -102;
    }

    // 3. Final check for zero length (handled in input_handler, but defensive check here)
    if (read_size == 0) {
        return 0; // Nothing to read
    }

    // Unsafe memory copy operation, which now reads past the end of 'buffer' if read_size > buffer_size.
    memcpy(destination, buffer, read_size);

    return (int)read_size; // Return the actual number of bytes read.
}