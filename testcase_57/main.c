#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "packet_parser.h"

// Utility function to print data in hex format
static void print_hex(const char *data, int len) {
    for (int i = 0; i < len; i++) {
        printf("%02X ", (unsigned char)data[i]);
    }
    printf("\n");
}

int main() {
    // 1. Setup simulated packet buffer (Source Buffer)
    // This buffer is guaranteed to be MAX_PACKET_SIZE (1024 bytes).
    char packet_buffer[MAX_PACKET_SIZE];
    
    // Initialize buffer with simulated network data (e.g., 0xAA followed by increasing sequence)
    for (size_t i = 0; i < MAX_PACKET_SIZE; i++) {
        packet_buffer[i] = (char)(0xAA + (i % 10));
    }

    printf("\n--- Packet Parser Utility ---\n");
    printf("Simulated packet buffer initialized (Size: %d bytes).\n", MAX_PACKET_SIZE);

    // 2. Get user input (Tainted Source)
    // The input handler ensures the value is positive and numeric.
    long requested_length_l = read_integer_input(
        "Enter desired data length to extract (Max allowed by destination: 512): "
    );

    if (requested_length_l <= 0) {
        fprintf(stderr, "Initialization failed due to invalid input.\n");
        return EXIT_FAILURE;
    }
    
    // Cast to size_t after validation. This value is the potential attack vector.
    size_t requested_length = (size_t)requested_length_l;

    // 3. Prepare destination buffer (Destination Buffer)
    char extracted_data[MAX_EXTRACT_SIZE];
    memset(extracted_data, 0, MAX_EXTRACT_SIZE);

    printf("Attempting to extract %zu bytes...\n", requested_length);

    // 4. Call the secure parser function
    // The security check against CWE-125 happens inside extract_data_segment,
    // comparing requested_length against MAX_PACKET_SIZE.
    int bytes_read = extract_data_segment(
        packet_buffer, 
        MAX_PACKET_SIZE, // Actual size of the source buffer
        requested_length, // Tainted input length
        extracted_data
    );

    // 5. Handle results securely
    if (bytes_read > 0) {
        printf("\n[SUCCESS] Successfully extracted %d bytes of data.\n", bytes_read);
        printf("Extracted Data (First 20 bytes): ");
        print_hex(extracted_data, (bytes_read > 20 ? 20 : bytes_read));
    } else if (bytes_read == 0) {
        printf("No data extracted (zero length requested).\n");
    } else {
        fprintf(stderr, "\n[FAILURE] Data extraction failed. Error code: %d. Check logs for security alerts.\n", bytes_read);
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}