#ifndef CONFIG_H
#define CONFIG_H

// Define the fixed size of the internal configuration buffer (e.g., 1024 bytes)
#define MAX_BUFFER_SIZE 1024

// Define the maximum allowed size for a single segment extraction request.
// This size also dictates the size of the result buffer in data_processor.h.
#define MAX_SEGMENT_SIZE 256

// Define the maximum valid index (0 to MAX_BUFFER_SIZE - 1)
#define MAX_VALID_INDEX (MAX_BUFFER_SIZE - 1)

#endif // CONFIG_H