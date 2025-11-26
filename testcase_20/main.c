#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stddef.h>

#include "data_processor.h"
#include "constants.h"

// Function to simulate reading data from an untrusted stream (e.g., network socket)
static size_t simulate_read_input(char *buffer, size_t max_buffer_size, size_t requested_length) {
    // In a real application, this would be the actual I/O call (e.g., recv, read).
    
    // Defensive check in the I/O layer (good practice, but not the primary defense)
    size_t bytes_to_read = requested_length;
    if (bytes_to_read > max_buffer_size) {
        bytes_to_read = max_buffer_size;
    }

    // Ensure the simulated data fits the expected structure
    if (bytes_to_read < EXPECTED_RECORD_SIZE) {
        return 0;
    }

    // Zero out the buffer before use (defense against information leakage)
    memset(buffer, 0, max_buffer_size);

    // Fill the buffer with dummy data
    // Data layout: ID (8 bytes), Name (64 bytes), Value (128 bytes)
    
    strncpy(buffer, "REC00001", FIELD_ID_SIZE);
    strncpy(buffer + FIELD_ID_SIZE, "TestName", FIELD_NAME_SIZE);
    strncpy(buffer + FIELD_ID_SIZE + FIELD_NAME_SIZE, "TestDataValue", FIELD_VALUE_SIZE);
    
    return bytes_to_read;
}

int main() {
    // Define the local, fixed-size buffer.
    char raw_buffer[MAX_BUFFER_SIZE];
    ParsedRecord record_output = {0};
    
    // --- Scenario 1: Secure, Valid Input ---
    size_t untrusted_length_1 = EXPECTED_RECORD_SIZE; 
    size_t actual_read_1 = simulate_read_input(raw_buffer, MAX_BUFFER_SIZE, untrusted_length_1);
    
    printf("--- Scenario 1: Valid Input (%zu bytes claimed) ---\n", untrusted_length_1);
    if (actual_read_1 > 0) {
        // Pass the claimed length (untrusted_length_1) to the processor
        ProcessingResult res = process_record(raw_buffer, untrusted_length_1, &record_output);
        
        if (res == SUCCESS) {
            printf("Processing successful.\n");
            printf("  ID: %s\n", record_output.id);
            printf("  Name: %s\n", record_output.name);
            printf("  Value: %s\n", record_output.value);
        } else {
            printf("Processing failed with error code: %d\n", res);
        }
    }
    
    // --- Scenario 2: Attack Attempt (CWE-125 Out-of-bounds Read) ---
    // Attacker provides a length much larger than the buffer size (MAX_BUFFER_SIZE).
    size_t untrusted_length_2 = MAX_BUFFER_SIZE + 500; 
    
    printf("\n--- Scenario 2: Malicious Oversized Input Attempt (%zu bytes claimed) ---\n", untrusted_length_2);
    
    // Simulate reading the maximum possible data into the buffer (MAX_BUFFER_SIZE bytes)
    size_t actual_read_2 = simulate_read_input(raw_buffer, MAX_BUFFER_SIZE, MAX_BUFFER_SIZE); 
    
    // Pass the untrusted, oversized length to the processor. 
    // The defense in data_processor.c must catch this.
    ProcessingResult res_attack = process_record(raw_buffer, untrusted_length_2, &record_output);
    
    if (res_attack == ERROR_INVALID_SIZE) {
        printf("Security Defense Triggered: Input size validation prevented potential CWE-125 OOB Read.\n");
    } else if (res_attack == SUCCESS) {
        printf("Processing successful (ERROR: Security defense failed).\n");
    } else {
        printf("Processing failed with unexpected error code: %d\n", res_attack);
    }

    return 0;
}