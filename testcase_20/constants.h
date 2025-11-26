#ifndef CONSTANTS_H
#define CONSTANTS_H

#include <stddef.h>

// Maximum size of the input buffer we are willing to process.
// This is the critical defense against the OOB read vulnerability.
#define MAX_BUFFER_SIZE 1024

// Define the structure of the fixed-size record fields
#define FIELD_ID_SIZE 8
#define FIELD_NAME_SIZE 64
#define FIELD_VALUE_SIZE 128

// Total expected size of a valid record
#define EXPECTED_RECORD_SIZE (FIELD_ID_SIZE + FIELD_NAME_SIZE + FIELD_VALUE_SIZE)

// Error Codes
typedef enum {
    SUCCESS = 0,
    ERROR_NULL_BUFFER = 1,
    ERROR_INVALID_SIZE = 2, // Used when untrusted size exceeds MAX_BUFFER_SIZE
    ERROR_PROCESSING_FAILED = 3
} ProcessingResult;

#endif // CONSTANTS_H