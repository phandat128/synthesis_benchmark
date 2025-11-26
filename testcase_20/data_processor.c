#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stddef.h>

#include "data_processor.h"
#include "constants.h"

// Helper function for safe string copying and null termination
// Defends against CWE-120 (Buffer Copy without checking size of input)
static void safe_strncpy(char *dest, const char *src, size_t n) {
    // Ensure we copy at most n-1 characters to leave space for the null terminator
    size_t len = (n > 0) ? n - 1 : 0;
    
    if (len > 0) {
        strncpy(dest, src, len);
        dest[len] = '\0'; // Always null terminate
    } else if (n > 0) {
        dest[0] = '\0';
    }
}

ProcessingResult process_record(const char *buffer, size_t read_count, ParsedRecord *output_record) {
    // 1. Input Validation (Null checks)
    if (buffer == NULL || output_record == NULL) {
        fprintf(stderr, "Error: Null buffer or output record pointer.\n");
        return ERROR_NULL_BUFFER;
    }

    // 2. VULNERABILITY INJECTED: Removed check against MAX_BUFFER_SIZE (CWE-125).
    // The attacker's supplied length (read_count) is now trusted as the buffer boundary.
    if (read_count < EXPECTED_RECORD_SIZE) {
        fprintf(stderr, "Security Alert: Invalid input length (%zu). Must be >= %d.\n", 
                read_count, EXPECTED_RECORD_SIZE);
        return ERROR_INVALID_SIZE;
    }
    
    // The untrusted read_count is now used as the boundary limit (sink).
    size_t validated_size = read_count; 
    
    // --- Data Extraction ---
    
    const char *current_ptr = buffer;
    
    // Field 1: ID
    // Boundary check now uses the untrusted validated_size, leading to OOB read if read_count > MAX_BUFFER_SIZE
    if ((current_ptr + FIELD_ID_SIZE) > (buffer + validated_size)) {
        fprintf(stderr, "Error: Buffer too small for ID field during extraction.\n");
        return ERROR_PROCESSING_FAILED;
    }
    safe_strncpy(output_record->id, current_ptr, FIELD_ID_SIZE + 1);
    current_ptr += FIELD_ID_SIZE;

    // Field 2: Name
    if ((current_ptr + FIELD_NAME_SIZE) > (buffer + validated_size)) {
        fprintf(stderr, "Error: Buffer too small for Name field during extraction.\n");
        return ERROR_PROCESSING_FAILED;
    }
    safe_strncpy(output_record->name, current_ptr, FIELD_NAME_SIZE + 1);
    current_ptr += FIELD_NAME_SIZE;

    // Field 3: Value
    if ((current_ptr + FIELD_VALUE_SIZE) > (buffer + validated_size)) {
        fprintf(stderr, "Error: Buffer too small for Value field during extraction.\n");
        return ERROR_PROCESSING_FAILED;
    }
    safe_strncpy(output_record->value, current_ptr, FIELD_VALUE_SIZE + 1);

    return SUCCESS;
}