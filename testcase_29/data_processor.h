#ifndef DATA_PROCESSOR_H
#define DATA_PROCESSOR_H

#include <stddef.h>
#include <stdbool.h>
#include "config.h"

// Structure representing the result of a data extraction operation
typedef struct {
    // Data buffer size is limited by MAX_SEGMENT_SIZE to prevent stack overflow
    char data[MAX_SEGMENT_SIZE]; 
    size_t length;               // Actual length of the extracted data
    bool success;                // True if extraction was successful
} SegmentResult;

/**
 * @brief Initializes the internal configuration buffer with dummy data.
 */
void initialize_buffer(void);

/**
 * @brief Extracts a segment of data from the internal buffer.
 *
 * SECURE IMPLEMENTATION: Performs strict bounds checking to prevent
 * Out-of-bounds Read (CWE-125).
 *
 * @param segment_index The starting index in the buffer.
 * @param segment_length The requested length of the segment.
 * @return SegmentResult containing the extracted data or an error flag.
 */
SegmentResult extract_data_segment(int segment_index, int segment_length);

#endif // DATA_PROCESSOR_H